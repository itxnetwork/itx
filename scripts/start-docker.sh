#!/bin/bash

KEY="dev0"
CHAINID="itx_9000-1"
MONIKER="mymoniker"
DATA_DIR=$(mktemp -d -t itx-datadir.XXXXX)

echo "create and add new keys"
./itxd keys add $KEY --home $DATA_DIR --no-backup --chain-id $CHAINID --algo "eth_secp256k1" --keyring-backend test
echo "init Itx with moniker=$MONIKER and chain-id=$CHAINID"
./itxd init $MONIKER --chain-id $CHAINID --home $DATA_DIR
echo "prepare genesis: Allocate genesis accounts"
./itxd add-genesis-account \
"$(./itxd keys show $KEY -a --home $DATA_DIR --keyring-backend test)" 1000000000000000000uitx,1000000000000000000stake \
--home $DATA_DIR --keyring-backend test
echo "prepare genesis: Sign genesis transaction"
./itxd gentx $KEY 1000000000000000000stake --keyring-backend test --home $DATA_DIR --keyring-backend test --chain-id $CHAINID
echo "prepare genesis: Collect genesis tx"
./itxd collect-gentxs --home $DATA_DIR
echo "prepare genesis: Run validate-genesis to ensure everything worked and that the genesis file is setup correctly"
./itxd validate-genesis --home $DATA_DIR

echo "starting itx node $i in background ..."
./itxd start --pruning=nothing --rpc.unsafe \
--keyring-backend test --home $DATA_DIR \
>$DATA_DIR/node.log 2>&1 & disown

echo "started itx node"
tail -f /dev/null