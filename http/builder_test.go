package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
	consensus "github.com/umbracle/go-eth-consensus"
)

func TestBuilderEndpoint(t *testing.T) {
	n := New("http://127.0.0.1:4011").Builder()

	t.Run("RegisterValidator", func(t *testing.T) {
		obj := []*SignedValidatorRegistration{
			{Message: &RegisterValidatorRequest{}},
		}
		err := n.RegisterValidator(obj)
		assert.NoError(t, err)
	})

	t.Run("GetExecutionPayload", func(t *testing.T) {
		_, err := n.GetExecutionPayload(1, [32]byte{}, [48]byte{})
		assert.NoError(t, err)
	})

	t.Run("SubmitBlindedBlock", func(t *testing.T) {
		obj := &consensus.SignedBlindedBeaconBlock{}
		_, err := n.SubmitBlindedBlock(obj)
		assert.NoError(t, err)
	})
}
