package spec

import consensus "github.com/umbracle/go-eth-consensus"

var Spec = &consensus.Spec{
	SecondsPerSlot:                   12,
	SlotsPerEpoch:                    32,
	MaxCommitteesPerSlot:             64,
	EpochsPerHistoricalVector:        65536,
	MinSeedLookAhead:                 1,
	MinEpochsToInactivityPenalty:     4,
	EffectiveBalanceIncrement:        1000000000,
	MaxEffectiveBalance:              32000000000,
	BaseRewardFactor:                 64,
	BaseRewardsPerEpoch:              4,
	TargetCommitteeSize:              128,
	ShardCommiteePeriod:              256,
	MaxSeedLookAhead:                 4,
	ChurnLimitQuotient:               65536,
	MinPerEpochChurnLimit:            4,
	EpochsPerSlashingsVector:         8192,
	MinSlashingPenaltyQuotient:       128,
	WhistleblowerRewardQuotient:      512,
	ProposerRewardQuotient:           8,
	MinValidatorWithdrawabilityDelay: 256,
	MinAttestationInclusionDelay:     1,
	EpochsPerEth1VotingPeriod:        64,
	SlotsPerHistoricalRoot:           8192,
	ProportionalSlashingsMultiplier:  1,
	HysteresisQuotient:               4,
	HysteresisDownwardMultiplier:     1,
	HysteresisUpwardMultiplier:       5,
	EjectionBalance:                  16000000000, // Gwei(2**4 * 10**9)
	InactivityPenaltyQuotient:        67108864,    // Gwei(2**26)
}
