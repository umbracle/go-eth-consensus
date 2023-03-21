package spec

import (
	"reflect"
	"testing"

	consensus "github.com/umbracle/go-eth-consensus"
)

func TestOpAttestation(t *testing.T) {
	type attestationTest struct {
		Attestation consensus.Attestation
		Pre         consensus.BeaconStatePhase0
		Post        consensus.BeaconStatePhase0
	}

	listTestData(t, "mainnet/phase0/operations/attestation/*/*", func(th *testHandler) {
		attestationTest := &attestationTest{}
		th.decodeFile("attestation", &attestationTest.Attestation)
		th.decodeFile("pre", &attestationTest.Pre)
		ok := th.decodeFile("post", &attestationTest.Post, true)

		if err := ProcessAttestation(&attestationTest.Pre, &attestationTest.Attestation); err != nil {
			if ok {
				t.Fatal(err)
			}
			return
		}

		if !ok {
			t.Fatal("it should fail")
		}
		if !reflect.DeepEqual(attestationTest.Pre, attestationTest.Post) {
			t.Fatal("bad")
		}
	})
}

func TestOpProcessAttesterSlashing(t *testing.T) {
	type processAttesterSlashingTest struct {
		Pre              consensus.BeaconStatePhase0
		Post             consensus.BeaconStatePhase0
		AttesterSlashing consensus.AttesterSlashing
	}

	listTestData(t, "mainnet/phase0/operations/attester_slashing/*/*", func(th *testHandler) {
		slashTest := &processAttesterSlashingTest{}
		th.decodeFile("pre", &slashTest.Pre)
		ok := th.decodeFile("post", &slashTest.Post, true)
		th.decodeFile("attester_slashing", &slashTest.AttesterSlashing)

		if err := ProcessAttesterSlashing(&slashTest.Pre, &slashTest.AttesterSlashing); err != nil {
			if ok {
				t.Fatal(err)
			}
			return
		}

		if !ok {
			t.Fatal("it should fail")
		}
		if !reflect.DeepEqual(slashTest.Pre, slashTest.Post) {
			t.Fatal("bad")
		}
	})
}

func TestOpProcessBlockBlockHeader(t *testing.T) {
	type blockHeaderTest struct {
		Pre   consensus.BeaconStatePhase0
		Post  consensus.BeaconStatePhase0
		Block consensus.BeaconBlockPhase0
	}

	listTestData(t, "mainnet/phase0/operations/block_header/*/*", func(th *testHandler) {
		blockHeaderTest := &blockHeaderTest{}
		th.decodeFile("block", &blockHeaderTest.Block)
		th.decodeFile("pre", &blockHeaderTest.Pre)
		ok := th.decodeFile("post", &blockHeaderTest.Post, true)

		if err := ProcessBlockHeader(&blockHeaderTest.Pre, &blockHeaderTest.Block); err != nil {
			if ok {
				t.Fatal(err)
			}
			return
		}

		if !ok {
			t.Fatal("it should fail")
		}
		if !reflect.DeepEqual(blockHeaderTest.Pre, blockHeaderTest.Post) {
			t.Fatal("bad")
		}
	})

}

func TestOpDeposit(t *testing.T) {
	type depositTest struct {
		Deposit consensus.Deposit
		Pre     consensus.BeaconStatePhase0
		Post    consensus.BeaconStatePhase0
	}

	listTestData(t, "mainnet/phase0/operations/deposit/*/*", func(th *testHandler) {
		depositTest := &depositTest{}
		th.decodeFile("deposit", &depositTest.Deposit)
		th.decodeFile("pre", &depositTest.Pre)
		ok := th.decodeFile("post", &depositTest.Post, true)

		if err := ProcessDeposit(&depositTest.Pre, &depositTest.Deposit); err != nil {
			if ok {
				t.Fatal(err)
			}
			return
		}

		if !ok {
			t.Fatal("it should fail")
		}
		if !reflect.DeepEqual(depositTest.Pre, depositTest.Post) {
			t.Fatal("bad")
		}
	})
}

func TestOpProposerSlashing(t *testing.T) {
	type proposerSlashingTest struct {
		ProposerSlashing consensus.ProposerSlashing
		Pre              consensus.BeaconStatePhase0
		Post             consensus.BeaconStatePhase0
	}

	listTestData(t, "mainnet/phase0/operations/proposer_slashing/*/*", func(th *testHandler) {
		proposerSlashingTest := &proposerSlashingTest{}
		th.decodeFile("proposer_slashing", &proposerSlashingTest.ProposerSlashing)
		th.decodeFile("pre", &proposerSlashingTest.Pre)
		ok := th.decodeFile("post", &proposerSlashingTest.Post, true)

		if err := ProcessProposerSlashing(&proposerSlashingTest.Pre, &proposerSlashingTest.ProposerSlashing); err != nil {
			if ok {
				t.Fatal(err)
			}
			return
		}

		if !ok {
			t.Fatal("it should fail")
		}
		if !reflect.DeepEqual(proposerSlashingTest.Pre, proposerSlashingTest.Post) {
			t.Fatal("bad")
		}
	})
}

func TestOpVoluntaryExit(t *testing.T) {
	type voluntaryExitTest struct {
		VoluntaryExit consensus.SignedVoluntaryExit
		Pre           consensus.BeaconStatePhase0
		Post          consensus.BeaconStatePhase0
	}

	listTestData(t, "mainnet/phase0/operations/voluntary_exit/*/*", func(th *testHandler) {
		voluntaryExitTest := &voluntaryExitTest{}
		th.decodeFile("voluntary_exit", &voluntaryExitTest.VoluntaryExit)
		th.decodeFile("pre", &voluntaryExitTest.Pre)
		ok := th.decodeFile("post", &voluntaryExitTest.Post, true)

		if err := ProcessVoluntaryExit(&voluntaryExitTest.Pre, &voluntaryExitTest.VoluntaryExit); err != nil {
			if ok {
				t.Fatal(err)
			}
			return
		}

		if !ok {
			t.Fatal("it should fail")
		}
		if !reflect.DeepEqual(voluntaryExitTest.Pre, voluntaryExitTest.Post) {
			t.Fatal("bad")
		}
	})
}
