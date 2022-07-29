#!/bin/bash

# Start beacon open api mock
docker run --init -p 4010:4010 stoplight/prism:4 mock -h 0.0.0.0 https://github.com/ethereum/beacon-APIs/releases/download/v2.3.0/beacon-node-oapi.json
