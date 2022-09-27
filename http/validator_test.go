package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
	consensus "github.com/umbracle/go-eth-consensus"
)

func TestValidatorEndpoint(t *testing.T) {
	n := New("http://127.0.0.1:4010", WithUntrackedKeys()).Validator()

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

	t.Run("GetBlock", func(t *testing.T) {
		t.Skip("due to graffiti")

		block := consensus.BeaconBlockAltair{}
		err := n.GetBlock(&block, 1, [96]byte{})
		assert.NoError(t, err)
	})

	t.Run("RequestAttestationData", func(t *testing.T) {
		_, err := n.RequestAttestationData(1, 1)
		assert.NoError(t, err)
	})

	t.Run("AggregateAttestation", func(t *testing.T) {
		_, err := n.AggregateAttestation(1, [32]byte{})
		assert.NoError(t, err)
	})

	t.Run("SyncCommitteeContribution", func(t *testing.T) {
		_, err := n.SyncCommitteeContribution(1, 1, [32]byte{})
		assert.NoError(t, err)
	})

	t.Run("BeaconCommitteeSubscriptions", func(t *testing.T) {
		err := n.BeaconCommitteeSubscriptions([]*BeaconCommitteeSubscription{{}})
		assert.NoError(t, err)
	})

	t.Run("SyncCommitteeSubscriptions", func(t *testing.T) {
		err := n.SyncCommitteeSubscriptions([]*SyncCommitteeSubscription{{}})
		assert.NoError(t, err)
	})

	t.Run("PrepareBeaconProposer", func(t *testing.T) {
		err := n.PrepareBeaconProposer([]*ProposalPreparation{{}})
		assert.NoError(t, err)
	})
}
