package spec

import (
	"reflect"
	"testing"

	consensus "github.com/umbracle/go-eth-consensus"
)

type rewardFunc func(state *consensus.BeaconStatePhase0) ([]uint64, []uint64)

func TestRewards(t *testing.T) {
	listTestData(t, "mainnet/phase0/rewards/basic/pyspec_tests/*", func(th *testHandler) {
		test := &specRewardTest{}
		test.Decode(th)

		cases := []struct {
			name  string
			fn    rewardFunc
			delta Deltas
		}{
			{"source", getSourceDeltas, test.SourceDeltas},
			{"target", getTargetDeltas, test.TargetDeltas},
			{"head", getHeadDeltas, test.HeadDeltas},
			{"inactivity", getInactivityPenaltyDeltas, test.InactivityPenaltyDeltas},
		}

		for _, c := range cases {
			rewards, penalties := c.fn(&test.Pre)

			if !reflect.DeepEqual(rewards, c.delta.Rewards) {
				t.Fatalf("bad '%s' rewards: %s", c.name, th.path)
			}
			if !reflect.DeepEqual(penalties, c.delta.Penalties) {
				t.Fatalf("bad '%s' penalties: %s", c.name, th.path)
			}
		}
	})
}

type specRewardTest struct {
	Pre                     consensus.BeaconStatePhase0
	HeadDeltas              Deltas
	InactivityPenaltyDeltas Deltas
	SourceDeltas            Deltas
	TargetDeltas            Deltas
}

func (s *specRewardTest) Decode(th *testHandler) {
	th.decodeFile("pre", &s.Pre)
	th.decodeFile("head_deltas", &s.HeadDeltas)
	th.decodeFile("inactivity_penalty_deltas", &s.InactivityPenaltyDeltas)
	th.decodeFile("source_deltas", &s.SourceDeltas)
	th.decodeFile("target_deltas", &s.TargetDeltas)
}
