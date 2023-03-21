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

	MaxCommitteesPerSlot uint64 `json:"MAX_COMMITTEES_PER_SLOT"`

	EpochsPerHistoricalVector uint64 `json:"EPOCHS_PER_HISTORICAL_VECTOR"`
	MinSeedLookAhead          uint64 `json:"MIN_SEED_LOOKAHEAD"`

	MinEpochsToInactivityPenalty uint64 `json:"MIN_EPOCHS_TO_INACTIVITY_PENALTY"`
	EffectiveBalanceIncrement    uint64 `json:"EFFECTIVE_BALANCE_INCREMENT"`

	BaseRewardFactor    uint64 `json:"BASE_REWARD_FACTOR"`
	BaseRewardsPerEpoch uint64 `json:"BASE_REWARDS_PER_EPOCH"`

	HysteresisQuotient           uint64 `json:"HYSTERESIS_QUOTIENT"`
	HysteresisDownwardMultiplier uint64 `json:"HYSTERESIS_DOWNWARD_MULTIPLIER"`
	HysteresisUpwardMultiplier   uint64 `json:"HYSTERESIS_UPWARD_MULTIPLIER"`

	EjectionBalance uint64 `json:"EJECTION_BALANCE"`

	TargetCommitteeSize uint64 `json:"TARGET_COMMITTEE_SIZE"`

	MaxEffectiveBalance uint64 `json:"MAX_EFFECTIVE_BALANCE"`

	EpochsPerEth1VotingPeriod uint64 `json:"EPOCHS_PER_ETH1_VOTING_PERIOD"`

	ProportionalSlashingsMultiplier uint64 `json:"PROPORTIONAL_SLASHING_MULTIPLIER"`
	SlotsPerHistoricalRoot          uint64 `json:"SLOTS_PER_HISTORICAL_ROOT"`
	SyncCommitteeSize               uint64 `json:"SYNC_COMMITTEE_SIZE"`
	MaxSeedLookAhead                uint64 `json:"MAX_SEED_LOOKAHEAD"`
	ShardCommiteePeriod             uint64 `json:"SHARD_COMMITTEE_PERIOD"`
	WhistleblowerRewardQuotient     uint64 `json:"WHISTLEBLOWER_REWARD_QUOTIENT"`
	ProposerRewardQuotient          uint64 `json:"PROPOSER_REWARD_QUOTIENT"`
	MinSlashingPenaltyQuotient      uint64 `json:"MIN_SLASHING_PENALTY_QUOTIENT"`

	InactivityPenaltyQuotient        uint64 `json:"INACTIVITY_PENALTY_QUOTIENT"`
	MinAttestationInclusionDelay     uint64 `json:"MIN_ATTESTATION_INCLUSION_DELAY"`
	MinValidatorWithdrawabilityDelay uint64 `json:"MIN_VALIDATOR_WITHDRAWABILITY_DELAY"`
	EpochsPerSlashingsVector         uint64 `json:"EPOCHS_PER_SLASHINGS_VECTOR"`

	ChurnLimitQuotient    uint64 `json:"CHURN_LIMIT_QUOTIENT"`
	MinPerEpochChurnLimit uint64 `json:"MIN_PER_EPOCH_CHURN_LIMIT"`

	// TargetAggregatorsPerCommittee defines the number of aggregators inside one committee.
	TargetAggregatorsPerCommittee uint64 `json:"TARGET_AGGREGATORS_PER_COMMITTEE"`

	// GenesisForkVersion is used to track fork version between state transitions.
	GenesisForkVersion Domain `json:"GENESIS_FORK_VERSION"`

	AltairForkVersion Domain `json:"ALTAIR_FORK_VERSION"`
	AltairForkEpoch   uint64 `json:"ALTAIR_FORK_EPOCH"`

	BellatrixForkVersion Domain `json:"BELLATRIX_FORK_VERSION"`
	BellatrixForkEpoch   uint64 `json:"BELLATRIX_FORK_EPOCH"`
}
