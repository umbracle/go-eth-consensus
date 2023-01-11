package spec

import (
	"testing"

	consensus "github.com/umbracle/go-eth-consensus"
)

func TestDeposit(t *testing.T) {
	th := &testHandler{
		t:    t,
		path: "../eth2.0-spec-tests/tests/mainnet/phase0/operations/deposit/pyspec_tests/new_deposit_under_max",
	}

	depositTest := &depositTest{}
	depositTest.Decode(th)
}

type depositTest struct {
	Deposit consensus.Deposit
	Pre     consensus.BeaconStatePhase0
	Post    consensus.BeaconStatePhase0
}

func (d *depositTest) Decode(th *testHandler) {
	th.decodeFile("deposit", &d.Deposit)
	th.decodeFile("pre", &d.Pre)
	th.decodeFile("post", &d.Post)
}
