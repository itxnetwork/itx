// Copyright 2022 Itx Foundation
// This file is part of the Itx Network packages.
//
// Itx is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The Itx packages are distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the Itx packages. If not, see https://github.com/itxnetwork/itx/blob/main/LICENSE

package app

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	"github.com/itxnetwork/itx/v12/utils"
)

// ScheduleForkUpgrade executes any necessary fork logic for based upon the current
// block height and chain ID (mainnet or testnet). It sets an upgrade plan once
// the chain reaches the pre-defined upgrade height.
//
// CONTRACT: for this logic to work properly it is required to:
//
//  1. Release a non-breaking patch version so that the chain can set the scheduled upgrade plan at upgrade-height.
//  2. Release the software defined in the upgrade-info
func (app *Itx) ScheduleForkUpgrade(ctx sdk.Context) {
	// NOTE: there are no testnet forks for the existing versions
	if !utils.IsMainnet(ctx.ChainID()) {
		return
	}

	upgradePlan := upgradetypes.Plan{
		Height: ctx.BlockHeight(),
	}

	// handle mainnet forks with their corresponding upgrade name and info
	switch ctx.BlockHeight() {
	default:
		// No-op
		return
	}

	// schedule the upgrade plan to the current block hight, effectively performing
	// a hard fork that uses the upgrade handler to manage the migration.
	if err := app.UpgradeKeeper.ScheduleUpgrade(ctx, upgradePlan); err != nil {
		panic(
			fmt.Errorf(
				"failed to schedule upgrade %s during BeginBlock at height %d: %w",
				upgradePlan.Name, ctx.BlockHeight(), err,
			),
		)
	}
}
