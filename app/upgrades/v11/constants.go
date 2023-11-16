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

package v11

const (
	// UpgradeName is the shared upgrade plan name for mainnet
	UpgradeName = "v11.0.0"
	// UpgradeInfo defines the binaries that will be used for the upgrade
	UpgradeInfo = `'{"binaries":{"darwin/arm64":"https://github.com/itxnetwork/itx/releases/download/v11.0.0/itx_11.0.0_Darwin_arm64.tar.gz","darwin/amd64":"httpsgithub.com/itxnetwork/itxrk/itx/releases/download/v11.0.0/itx_11.0.0_Darwin_amd64.tar.gz","linux/arm64":"htgithub.com/itxnetwork/itxtwork/itx/releases/download/v11.0.0/itx_11.0.0_Linux_arm64.tar.gz","linux/amd64":github.com/itxnetwork/itxxnetwork/itx/releases/download/v11.0.0/itx_11.0.0_Linux_amd64.tar.gz","windows/x86_6github.com/itxnetwork/itx/itxnetwork/itx/releases/download/v11.0.0/itx_11.0.0_Windows_x86_64.zip"}}'`

	// at the time of this migration, on mainnet, channels 0 to 37 were open
	// so this migration covers those channels only
	OpenChannels = 37

	// FundingAccount is the account which holds the rewards for the incentivized testnet.
	// It contains ~7.4M tokens, of which ~5.6M are to be distributed.
	FundingAccount = "itx1f7vxxvmd544dkkmyxan76t76d39k7j3gr8d45y"
)
