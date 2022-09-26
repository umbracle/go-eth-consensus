#!/bin/bash

# Start beacon open api mock
docker run --init --name prism-beacon -d -p 4010:4010 stoplight/prism:4.10.3 mock -h 0.0.0.0 https://github.com/ethereum/beacon-APIs/releases/download/v2.3.0/beacon-node-oapi.json

# Wait for beacon mock api
grep -m 1 "Prism is listening" <(docker logs prism-beacon -f)

# Start builder api mock
docker run --init --name prism-builder -d -p 4011:4010 stoplight/prism:4.10.3 mock -h 0.0.0.0 https://github.com/ethereum/builder-specs/releases/download/v0.2.0/builder-oapi.yaml

# Wait for builder mock api
grep -m 1 "Prism is listening" <(docker logs prism-builder -f)
