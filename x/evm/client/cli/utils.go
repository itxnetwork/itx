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
package cli

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/ethereum/go-ethereum/common"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func accountToHex(addr string) (string, error) {
	if strings.HasPrefix(addr, sdk.GetConfig().GetBech32AccountAddrPrefix()) {
		// Check to see if address is Cosmos bech32 formatted
		toAddr, err := sdk.AccAddressFromBech32(addr)
		if err != nil {
			return "", errors.Wrap(err, "must provide a valid Bech32 address")
		}
		ethAddr := common.BytesToAddress(toAddr.Bytes())
		return ethAddr.Hex(), nil
	}

	if !strings.HasPrefix(addr, "0x") {
		addr = "0x" + addr
	}

	valid := common.IsHexAddress(addr)
	if !valid {
		return "", fmt.Errorf("%s is not a valid Ethereum or Cosmos address", addr)
	}

	ethAddr := common.HexToAddress(addr)

	return ethAddr.Hex(), nil
}

func formatKeyToHash(key string) string {
	if !strings.HasPrefix(key, "0x") {
		key = "0x" + key
	}

	ethkey := common.HexToHash(key)

	return ethkey.Hex()
}
