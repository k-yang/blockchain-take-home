#!/usr/bin/env bash

rm -rf $HOME/.blog
BLOGD_BIN=$(which blogd)
if [ -z "$BLOGD_BIN" ]; then
    GOBIN=$(go env GOPATH)/bin
    BLOGD_BIN=$(which $GOBIN/blogd)
fi

if [ -z "$BLOGD_BIN" ]; then
    echo "please verify blogd is installed"
    exit 1
fi

# configure blogd
$BLOGD_BIN config set client chain-id blog
$BLOGD_BIN config set client keyring-backend test
$BLOGD_BIN keys add alice
$BLOGD_BIN keys add bob
$BLOGD_BIN init test --chain-id blog --default-denom stake
# update genesis
$BLOGD_BIN genesis add-genesis-account alice 10000000stake --keyring-backend test
$BLOGD_BIN genesis add-genesis-account bob 1000stake --keyring-backend test
$BLOGD_BIN genesis add-genesis-account cosmos1vq7mchr8585n0a4q75c9fjr9qpfygvxz8fv2vc 3000stake --keyring-backend test

# create default validator
$BLOGD_BIN genesis gentx alice 1000000stake --chain-id blog
$BLOGD_BIN genesis collect-gentxs