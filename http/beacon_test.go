package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
	consensus "github.com/umbracle/go-eth-consensus"
)

func TestBeaconEndpoint(t *testing.T) {
	n := New("http://127.0.0.1:4010", WithUntrackedKeys()).Beacon()

	t.Run("Genesis", func(t *testing.T) {
		_, err := n.Genesis()
		assert.NoError(t, err)
	})

	t.Run("SubmitCommitteeDuties", func(t *testing.T) {
		err := n.SubmitCommitteeDuties([]*consensus.SyncCommitteeMessage{})
		assert.NoError(t, err)
	})

	t.Run("GetRoot", func(t *testing.T) {
		_, err := n.GetRoot(Finalized)
		assert.NoError(t, err)
	})

	t.Run("GetFork", func(t *testing.T) {
		_, err := n.GetFork(Finalized)
		assert.NoError(t, err)
	})

	t.Run("GetFinalityCheckpoints", func(t *testing.T) {
		_, err := n.GetFinalityCheckpoints(Finalized)
		assert.NoError(t, err)
	})

	t.Run("GetValidators", func(t *testing.T) {
		_, err := n.GetValidators(Finalized)
		assert.NoError(t, err)
	})

	t.Run("GetValidatorByPubKey", func(t *testing.T) {
		_, err := n.GetValidatorByPubKey("0x1", Slot(1))
		assert.NoError(t, err)
	})

	t.Run("PublishAttestations", func(t *testing.T) {
		err := n.PublishAttestations([]*consensus.Attestation{})
		assert.NoError(t, err)
	})

	t.Run("GetBlock", func(t *testing.T) {
		t.Skip("graffiti TODO")

		var out consensus.BeaconBlockPhase0
		_, err := n.GetBlock(Slot(1), &out)
		assert.NoError(t, err)
	})

	t.Run("GetBlockHeader", func(t *testing.T) {
		_, err := n.GetBlockHeader(Finalized)
		assert.NoError(t, err)
	})

	t.Run("GetBlockRoot", func(t *testing.T) {
		_, err := n.GetBlockRoot(Head)
		assert.NoError(t, err)
	})

	t.Run("GetBlockAttestations", func(t *testing.T) {
		_, err := n.GetBlockAttestations(Genesis)
		assert.NoError(t, err)
	})
}
