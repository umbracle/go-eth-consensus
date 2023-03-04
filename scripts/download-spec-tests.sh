#!/bin/bash

VERSION=$1
REPO_NAME=eth2.0-spec-tests

# Remove dir if it already exists
rm -rf $REPO_NAME
mkdir $REPO_NAME

function download {
    OUTPUT=$1.tar.gz
    DOWNLOAD_URL=https://github.com/ethereum/$REPO_NAME/releases/download/$VERSION/$OUTPUT
    wget $DOWNLOAD_URL -O $OUTPUT
    tar -xzf $OUTPUT -C $REPO_NAME
    rm $OUTPUT
}

download "mainnet"
download "general"

# Download bls tests
mkdir $REPO_NAME/bls
wget https://github.com/ethereum/bls12-381-tests/releases/download/v0.1.2/bls_tests_json.tar.gz -O - | tar -xz -C eth2.0-spec-tests/bls
