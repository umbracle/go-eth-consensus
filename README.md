# Go-eth-consensus [![Godoc](https://godoc.org/github.com/umbracle/go-eth-consensus?status.svg)](https://godoc.org/github.com/umbracle/go-eth-consensus)

Go-eth-consensus is a suite of Go utilities to interact with the Ethereum consensus layer.

The core of this library was initially part of [eth2-validator](https://github.com/umbracle/eth2-validator). However, as [other](https://github.com/umbracle/viewpoint) projects started to mature, it became necessary to create a unified library to reduce code duplication and increase consistency.

## Features

**Consensus data types**. Full set of data types (up to Bellatrix) in `structs.go` at root. It includes the SSZ encoding for each one using [`fastssz`](https://github.com/ferranbt/fastssz). Each type is end-to-end tested with the official consensus spec tests.

**Http client**. Lightweight implementation for the [Beacon](https://ethereum.github.io/beacon-APIs) and [Builder](https://ethereum.github.io/builder-specs) OpenAPI spec. For usage and examples see the [Godoc](https://pkg.go.dev/github.com/umbracle/go-eth-consensus/http). The endpoints are tested against a real server that mocks the OpenAPI spec.

**Chaintime**. Simple utilities to interact with slot times and epochs.

**BLS**. Abstraction to sign, recover and store (with keystore format) BLS keys. It includes two implementations: [blst](https://github.com/supranational/blst) with cgo and [kilic/bls12-381](https://github.com/kilic/bls12-381) with pure Go. The build flag `CGO_ENABLED` determines which library is used.

## Installation

```
$ go get github.com/umbracle/go-eth-consensus
```
