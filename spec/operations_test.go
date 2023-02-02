package spec

import (
	"fmt"
	"reflect"
	"testing"

	consensus "github.com/umbracle/go-eth-consensus"
)

func TestAttestation(t *testing.T) {
	type attestationTest struct {
		Attestation consensus.Attestation
		Pre         consensus.BeaconStatePhase0
		Post        consensus.BeaconStatePhase0
	}

	listTestData(t, "mainnet/phase0/operations/attestation/*/*", func(th *testHandler) {
		attestationTest := &attestationTest{}
		th.decodeFile("attestation", &attestationTest.Attestation)
		th.decodeFile("pre", &attestationTest.Pre)
		th.decodeFile("post", &attestationTest.Post, true)

		fmt.Println("x")
	})
}

func TestProcessAttesterSlashing(t *testing.T) {
	type processAttesterSlashingTest struct {
		Pre              consensus.BeaconStatePhase0
		Post             consensus.BeaconStatePhase0
		AttesterSlashing consensus.AttesterSlashing
	}

	listTestData(t, "mainnet/phase0/operations/attester_slashing/*/*", func(th *testHandler) {
		slashTest := &processAttesterSlashingTest{}
		th.decodeFile("pre", &slashTest.Pre)
		th.decodeFile("post", &slashTest.Post, true)
		th.decodeFile("attester_slashing", &slashTest.AttesterSlashing)

		ProcessAttesterSlashing(&slashTest.Pre, &slashTest.AttesterSlashing)
	})
}

func TestProcessBlockBlockHeader(t *testing.T) {
	type blockHeaderTest struct {
		Pre   consensus.BeaconStatePhase0
		Post  consensus.BeaconStatePhase0
		Block consensus.BeaconBlockPhase0
	}

	listTestData(t, "mainnet/phase0/operations/block_header/*/*", func(th *testHandler) {
		blockHeaderTest := &blockHeaderTest{}
		th.decodeFile("block", &blockHeaderTest.Block)
		th.decodeFile("pre", &blockHeaderTest.Pre)
		th.decodeFile("post", &blockHeaderTest.Post, true)

		ProcessBlockHeader(&blockHeaderTest.Pre, &blockHeaderTest.Block)
	})

}

func TestDeposit(t *testing.T) {
	type depositTest struct {
		Deposit consensus.Deposit
		Pre     consensus.BeaconStatePhase0
		Post    consensus.BeaconStatePhase0
	}

	listTestData(t, "mainnet/phase0/operations/deposit/*/*", func(th *testHandler) {
		depositTest := &depositTest{}
		th.decodeFile("deposit", &depositTest.Deposit)
		th.decodeFile("pre", &depositTest.Pre)
		th.decodeFile("post", &depositTest.Post, true)

		ProcessDeposit(&depositTest.Pre, &depositTest.Deposit)

		if !reflect.DeepEqual(depositTest.Pre, depositTest.Post) {
			t.Fatal("bad")
		}
	})
}

func TestProposerSlashing(t *testing.T) {
	type proposerSlashingTest struct {
		ProposerSlashing consensus.ProposerSlashing
		Pre              consensus.BeaconStatePhase0
		Post             consensus.BeaconStatePhase0
	}

	listTestData(t, "mainnet/phase0/operations/proposer_slashing/*/*", func(th *testHandler) {
		proposerSlashingTest := &proposerSlashingTest{}
		th.decodeFile("proposer_slashing", &proposerSlashingTest.ProposerSlashing)
		th.decodeFile("pre", &proposerSlashingTest.Pre)
		th.decodeFile("post", &proposerSlashingTest.Post, true)

		ProcessProposerSlashing(&proposerSlashingTest.Pre, &proposerSlashingTest.ProposerSlashing)
	})
}

func TestVoluntaryExit(t *testing.T) {
	type voluntaryExitTest struct {
		VoluntaryExit consensus.SignedVoluntaryExit
		Pre           consensus.BeaconStatePhase0
		Post          consensus.BeaconStatePhase0
	}

	listTestData(t, "mainnet/phase0/operations/voluntary_exit/*/*", func(th *testHandler) {
		voluntaryExitTest := &voluntaryExitTest{}
		th.decodeFile("voluntary_exit", &voluntaryExitTest.VoluntaryExit)
		th.decodeFile("pre", &voluntaryExitTest.Pre)
		th.decodeFile("post", &voluntaryExitTest.Post, true)

		ProcessVoluntaryExit(&voluntaryExitTest.Pre, &voluntaryExitTest.VoluntaryExit)
	})
}
