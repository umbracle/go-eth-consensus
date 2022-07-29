package consensus

type Spec struct {
	GenesisSlot  uint64 `json:"GENESIS_SLOT"`  // GenesisSlot represents the first canonical slot number of the beacon chain.
	GenesisEpoch uint64 `json:"GENESIS_EPOCH"` // GenesisEpoch represents the first canonical epoch number of the beacon chain.

	// Fork values
	GenesisForkVersion Domain `json:"GENESIS_FORK_VERSION"` // GenesisForkVersion is used to track fork version between state transitions.

	AltairForkVersion Domain `json:"ALTAIR_FORK_VERSION"`
	AltairForkEpoch   uint64 `json:"ALTAIR_FORK_EPOCH"`

	BellatrixForkVersion Domain `json:"BELLATRIX_FORK_VERSION"`
	BellatrixForkEpoch   uint64 `json:"BELLATRIX_FORK_EPOCH"`
}
