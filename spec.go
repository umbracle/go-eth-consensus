package consensus

type Spec struct {
	// GenesisSlot represents the first canonical slot number of the beacon chain.
	GenesisSlot uint64 `json:"GENESIS_SLOT"`

	// GenesisEpoch represents the first canonical epoch number of the beacon chain.
	GenesisEpoch uint64 `json:"GENESIS_EPOCH"`

	// SecondsPerSlot is how many seconds are in a single slot.
	SecondsPerSlot uint64 `json:"SECONDS_PER_SLOT"`

	// SlotsPerEpoch is the number of slots in an epoch.
	SlotsPerEpoch uint64 `json:"SLOTS_PER_EPOCH"`

	MaxCommitteesPerSlot uint64 // TODO

	EpocsPerHistoricalVector uint64
	MinSeedLookAhead         uint64

	MinEpochsToInactivityPenalty uint64
	EffectiveBalanceIncrement    uint64

	BaseRewardFactor    uint64
	BaseRewardsPerEpoch uint64

	TargetCommitteeSize uint64

	SyncCommitteeSize uint64 `json:"SYNC_COMMITTEE_SIZE"`

	// TargetAggregatorsPerCommittee defines the number of aggregators inside one committee.
	TargetAggregatorsPerCommittee uint64 `json:"TARGET_AGGREGATORS_PER_COMMITTEE"`

	// GenesisForkVersion is used to track fork version between state transitions.
	GenesisForkVersion Domain `json:"GENESIS_FORK_VERSION"`

	AltairForkVersion Domain `json:"ALTAIR_FORK_VERSION"`
	AltairForkEpoch   uint64 `json:"ALTAIR_FORK_EPOCH"`

	BellatrixForkVersion Domain `json:"BELLATRIX_FORK_VERSION"`
	BellatrixForkEpoch   uint64 `json:"BELLATRIX_FORK_EPOCH"`
}
