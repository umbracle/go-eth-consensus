package spec

import (
	"fmt"
	"testing"

	consensus "github.com/umbracle/go-eth-consensus"
)

/*
func TestOperationsX(t *testing.T) {
	th := &testHandler{
		path: "../eth2.0-spec-tests/tests/mainnet/phase0/operations/attestation/pyspec_tests/after_epoch_slots",
	}

	test := &opAttestationTest{
		Pre:         &consensus.BeaconStatePhase0{},
		Attestation: &consensus.Attestation{},
	}
	th.decodeFile(t, "attestation", test.Attestation)
	th.decodeFile(t, "pre", test.Pre)
}

type opAttestationTest struct {
	Pre         *consensus.BeaconStatePhase0
	Attestation *consensus.Attestation
}
*/

func TestRewards(t *testing.T) {
	th := &testHandler{
		t:    t,
		path: "../eth2.0-spec-tests/tests/mainnet/phase0/rewards/basic/pyspec_tests/full_all_correct",
	}
	/*
		listTestData(path, func(th *testHandler) {
			fmt.Println(th.path)
		})
	*/

	test := &specRewardTest{}
	test.Decode(th)

	fmt.Println(getSourceDeltas(&test.Pre))

	fmt.Println(test.SourceDeltas)
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
	th.decodeFile("target_deltas", &s.SourceDeltas)
}
