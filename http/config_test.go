package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigEndpoint(t *testing.T) {
	n := New("http://127.0.0.1:4010", WithUntrackedKeys()).Config()

	t.Run("ForkSchedule", func(t *testing.T) {
		_, err := n.ForkSchedule()
		assert.NoError(t, err)
	})

	t.Run("Spec", func(t *testing.T) {
		t.Skip("due to lack of 'data'")
		// the Spec example marshals without a 'data' field.
		// It has been addressed on master but it has not been released yet.

		_, err := n.Spec()
		assert.NoError(t, err)
	})

	t.Run("DepositContract", func(t *testing.T) {
		_, err := n.DepositContract()
		assert.NoError(t, err)
	})
}
