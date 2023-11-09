// Copyright 2022 Evmos Foundation
// This file is part of the Evmos Network packages.
//
// Evmos is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The Evmos packages are distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the Evmos packages. If not, see https://github.com/evmos/evmos/blob/main/LICENSE

package keeper

import (
	"github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/evmos/evmos/v12/x/inflation/types"
)

// EndBlocker checks if the airdrop claiming period has ended in order to
// process the clawback of unclaimed tokens
func (k Keeper) EndBlocker(ctx sdk.Context) {
	params := k.GetParams(ctx)
	logger := k.Logger(ctx)

	mintedCoin := sdk.Coin{
		Denom:  params.MintDenom,
		Amount: k.GetMintAmount(ctx, ctx.BlockHeight()),
	}

	logger.Info("Inflation", "mint amount", mintedCoin.Amount)
	logger.Info("Inflation", "mint denom", mintedCoin.Denom)

	staking, incentives, communityPool, err := k.MintAndAllocateInflation(ctx, mintedCoin, params)
	if err != nil {
		panic(err)
	}

	defer func() {
		stakingAmt := staking.AmountOfNoDenomValidation(mintedCoin.Denom)
		incentivesAmt := incentives.AmountOfNoDenomValidation(mintedCoin.Denom)
		cpAmt := communityPool.AmountOfNoDenomValidation(mintedCoin.Denom)

		if mintedCoin.Amount.IsInt64() {
			telemetry.IncrCounterWithLabels(
				[]string{types.ModuleName, "allocate", "total"},
				float32(mintedCoin.Amount.Int64()),
				[]metrics.Label{telemetry.NewLabel("denom", mintedCoin.Denom)},
			)
		}
		if stakingAmt.IsInt64() {
			telemetry.IncrCounterWithLabels(
				[]string{types.ModuleName, "allocate", "staking", "total"},
				float32(stakingAmt.Int64()),
				[]metrics.Label{telemetry.NewLabel("denom", mintedCoin.Denom)},
			)
		}
		if incentivesAmt.IsInt64() {
			telemetry.IncrCounterWithLabels(
				[]string{types.ModuleName, "allocate", "incentives", "total"},
				float32(incentivesAmt.Int64()),
				[]metrics.Label{telemetry.NewLabel("denom", mintedCoin.Denom)},
			)
		}
		if cpAmt.IsInt64() {
			telemetry.IncrCounterWithLabels(
				[]string{types.ModuleName, "allocate", "community_pool", "total"},
				float32(cpAmt.Int64()),
				[]metrics.Label{telemetry.NewLabel("denom", mintedCoin.Denom)},
			)
		}
	}()
}
