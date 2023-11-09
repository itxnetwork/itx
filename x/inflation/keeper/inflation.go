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
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	evmostypes "github.com/evmos/evmos/v12/types"

	utils "github.com/evmos/evmos/v12/utils"
	incentivestypes "github.com/evmos/evmos/v12/x/incentives/types"
	"github.com/evmos/evmos/v12/x/inflation/types"
)

// 200M token at year 4 allocated to the team
var teamAlloc = sdk.NewInt(200_000_000).Mul(evmostypes.PowerReduction)

// retrieve the fixed amount to mint per block.
func (k Keeper) GetMintAmount(ctx sdk.Context, currentBlockHeight int64) sdk.Int {

	// TODO read values from parameters

	// Retrieve the start block height from parameters
	startBlock := int64(1) // k.GetStartBlock()

	// number of blocks represent a four-year period
	// seconds in year / block time in seconds
	// (365.25 * 24 * 3600) / 5
	blocksPerHalving := int64(6311520)

	// Calculate the number of halvings that have occurred.
	halvings := (currentBlockHeight - startBlock) / blocksPerHalving

	// Start with the initial mint amount and halve it for each halving period that has passed.

	// Create a big.Int representing 10^18
	bigTenToTheEighteenth := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)

	// Initialize a new sdk.Int with 10^18
	tenToTheEighteenth := sdk.NewIntFromBigInt(bigTenToTheEighteenth)

	// TODO Initial mint amount to be set in the parameters.
	mintAmount := tenToTheEighteenth

	for i := int64(0); i < halvings; i++ {
		mintAmount = mintAmount.Quo(sdk.NewInt(2)) // Halve the mint amount.
		if mintAmount.IsZero() {
			// Once the mint amount reaches zero, stop minting new tokens.
			return sdk.ZeroInt()
		}
	}

	return mintAmount
}

// MintAndAllocateInflation performs inflation minting and allocation
func (k Keeper) MintAndAllocateInflation(
	ctx sdk.Context,
	coin sdk.Coin,
	params types.Params,
) (
	staking, incentives, communityPool sdk.Coins,
	err error,
) {
	// skip as no coins need to be minted
	if coin.Amount.IsNil() || !coin.Amount.IsPositive() {
		return nil, nil, nil, nil
	}

	// Mint coins for distribution
	if err := k.MintCoins(ctx, coin); err != nil {
		return nil, nil, nil, err
	}

	// Allocate minted coins according to allocation proportions (staking, usage
	// incentives, community pool)
	return k.AllocateExponentialInflation(ctx, coin, params)
}

// MintCoins implements an alias call to the underlying supply keeper's
// MintCoins to be used in BeginBlocker.
func (k Keeper) MintCoins(ctx sdk.Context, coin sdk.Coin) error {
	coins := sdk.Coins{coin}
	return k.bankKeeper.MintCoins(ctx, types.ModuleName, coins)
}

// AllocateExponentialInflation allocates coins from the inflation to external
// modules according to allocation proportions:
//   - staking rewards -> sdk `auth` module fee collector
//   - usage incentives -> `x/incentives` module
//   - community pool -> `sdk `distr` module community pool
func (k Keeper) AllocateExponentialInflation(
	ctx sdk.Context,
	mintedCoin sdk.Coin,
	params types.Params,
) (
	staking, incentives, communityPool sdk.Coins,
	err error,
) {
	distribution := params.InflationDistribution

	// Allocate staking rewards into fee collector account
	staking = sdk.Coins{k.GetProportions(ctx, mintedCoin, distribution.StakingRewards)}

	if err := k.bankKeeper.SendCoinsFromModuleToModule(
		ctx,
		types.ModuleName,
		k.feeCollectorName,
		staking,
	); err != nil {
		return nil, nil, nil, err
	}

	// Allocate usage incentives to incentives module account
	incentives = sdk.Coins{k.GetProportions(ctx, mintedCoin, distribution.UsageIncentives)}

	if err = k.bankKeeper.SendCoinsFromModuleToModule(
		ctx,
		types.ModuleName,
		incentivestypes.ModuleName,
		incentives,
	); err != nil {
		return nil, nil, nil, err
	}

	// Allocate community pool amount (remaining module balance) to community
	// pool address
	moduleAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)
	inflationBalance := k.bankKeeper.GetAllBalances(ctx, moduleAddr)

	err = k.distrKeeper.FundCommunityPool(
		ctx,
		inflationBalance,
		moduleAddr,
	)
	if err != nil {
		return nil, nil, nil, err
	}

	return staking, incentives, communityPool, nil
}

// GetAllocationProportion calculates the proportion of coins that is to be
// allocated during inflation for a given distribution.
func (k Keeper) GetProportions(
	_ sdk.Context,
	coin sdk.Coin,
	distribution sdk.Dec,
) sdk.Coin {
	return sdk.Coin{
		Denom:  coin.Denom,
		Amount: sdk.NewDecFromInt(coin.Amount).Mul(distribution).TruncateInt(),
	}
}

// BondedRatio the fraction of the staking tokens which are currently bonded
// It doesn't consider team allocation for inflation
func (k Keeper) BondedRatio(ctx sdk.Context) sdk.Dec {
	stakeSupply := k.stakingKeeper.StakingTokenSupply(ctx)

	isMainnet := utils.IsMainnet(ctx.ChainID())

	if !stakeSupply.IsPositive() || (isMainnet && stakeSupply.LTE(teamAlloc)) {
		return sdk.ZeroDec()
	}

	// don't count team allocation in bonded ratio's stake supple
	if isMainnet {
		stakeSupply = stakeSupply.Sub(teamAlloc)
	}

	return sdk.NewDecFromInt(k.stakingKeeper.TotalBondedTokens(ctx)).QuoInt(stakeSupply)
}

// GetCirculatingSupply returns the bank supply of the mintDenom excluding the
// team allocation in the first year
func (k Keeper) GetCirculatingSupply(ctx sdk.Context, mintDenom string) sdk.Dec {
	circulatingSupply := sdk.NewDecFromInt(k.bankKeeper.GetSupply(ctx, mintDenom).Amount)
	teamAllocation := sdk.NewDecFromInt(teamAlloc)
	// Consider team allocation only on mainnet chain id
	if utils.IsMainnet(ctx.ChainID()) {
		circulatingSupply = circulatingSupply.Sub(teamAllocation)
	}
	return circulatingSupply
}

// GetInflationRate returns the inflation rate for the current period.
func (k Keeper) GetInflationRate(ctx sdk.Context, mintDenom string) sdk.Dec {
	epp := k.GetEpochsPerPeriod(ctx)
	if epp == 0 {
		return sdk.ZeroDec()
	}
	epochMintProvision := k.GetEpochMintProvision(ctx)
	if epochMintProvision.IsZero() {
		return sdk.ZeroDec()
	}
	epochsPerPeriod := sdk.NewDec(epp)
	circulatingSupply := k.GetCirculatingSupply(ctx, mintDenom)
	if circulatingSupply.IsZero() {
		return sdk.ZeroDec()
	}
	// EpochMintProvision * 365 / circulatingSupply * 100
	return epochMintProvision.Mul(epochsPerPeriod).Quo(circulatingSupply).Mul(sdk.NewDec(100))
}

// GetEpochMintProvision retrieves necessary params KV storage
// and calculate EpochMintProvision
func (k Keeper) GetEpochMintProvision(ctx sdk.Context) sdk.Dec {
	return types.CalculateEpochMintProvision(
		k.GetParams(ctx),
		k.GetPeriod(ctx),
		k.GetEpochsPerPeriod(ctx),
		k.BondedRatio(ctx),
	)
}
