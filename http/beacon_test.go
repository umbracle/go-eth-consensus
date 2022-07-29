package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBeaconEndpoint(t *testing.T) {
	n := New("http://127.0.0.1:4010").Beacon()

	t.Run("Genesis", func(t *testing.T) {
		_, err := n.Genesis()
		assert.NoError(t, err)
	})

	t.Run("SubmitCommitteeDuties", func(t *testing.T) {
		err := n.SubmitCommitteeDuties([]*SyncCommitteeMessage{})
		assert.NoError(t, err)
	})

	t.Run("GetValidatorByPubKey", func(t *testing.T) {
		_, err := n.GetValidatorByPubKey("0x1")
		assert.NoError(t, err)
	})

	t.Run("GetHeadBlockRoot", func(t *testing.T) {
		_, err := n.GetHeadBlockRoot()
		assert.NoError(t, err)
	})
}
