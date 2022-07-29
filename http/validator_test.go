package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidatorEndpoint(t *testing.T) {
	n := New("http://127.0.0.1:4010").Validator()

	b32 := make([]byte, 32)
	//b96 := make([]byte, 96)

	t.Run("GetAttesterDuties", func(t *testing.T) {
		_, err := n.GetAttesterDuties(1, []string{"1"})
		assert.NoError(t, err)
	})

	t.Run("GetProposerDuties", func(t *testing.T) {
		_, err := n.GetProposerDuties(1)
		assert.NoError(t, err)
	})

	t.Run("GetCommitteeSyncDuties", func(t *testing.T) {
		_, err := n.GetCommitteeSyncDuties(1, []string{"1"})
		assert.NoError(t, err)
	})

	/*
		t.Run("GetBlock", func(t *testing.T) {
			_, err := n.GetBlock(1, randao)
			assert.NoError(t, err)
		})
	*/

	t.Run("RequestAttestationData", func(t *testing.T) {
		_, err := n.RequestAttestationData(1, 1)
		assert.NoError(t, err)
	})

	t.Run("AggregateAttestation", func(t *testing.T) {
		_, err := n.AggregateAttestation(1, b32)
		assert.NoError(t, err)
	})

	//t.Run("PublishAggregateAndProof", func(t *testing.T) {
	//	panic("TODO")
	//})

	t.Run("SyncCommitteeContribution", func(t *testing.T) {
		_, err := n.SyncCommitteeContribution(1, 1, b32)
		assert.NoError(t, err)
	})
}
