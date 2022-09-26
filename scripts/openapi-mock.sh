#!/bin/bash

# Start beacon open api mock
docker run --init -d -p 4010:4010 stoplight/prism:4 mock -h 0.0.0.0 https://github.com/ethereum/beacon-APIs/releases/download/v2.3.0/beacon-node-oapi.json

# Start builder api mock
docker run --init -d -p 4011:4010 stoplight/prism:4 mock -h 0.0.0.0 https://github.com/ethereum/builder-specs/releases/download/v0.2.0/builder-oapi.yaml

# wait for the endpoints to be available
sleep 5
