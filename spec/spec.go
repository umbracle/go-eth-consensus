package spec

import consensus "github.com/umbracle/go-eth-consensus"

var Spec = &consensus.Spec{
	SecondsPerSlot:               12,
	SlotsPerEpoch:                32,
	MaxCommitteesPerSlot:         64,
	EpocsPerHistoricalVector:     65536,
	MinSeedLookAhead:             1,
	MinEpochsToInactivityPenalty: 4,
	EffectiveBalanceIncrement:    1000000000,
	MaxEffectiveBalance:          32000000000,
	BaseRewardFactor:             64,
	BaseRewardsPerEpoch:          4,
	TargetCommitteeSize:          128,
	ShardCommiteePeriod:          256,
}
