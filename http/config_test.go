package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigEndpoint(t *testing.T) {
	n := New("http://127.0.0.1:4010", WithUntrackedKeys()).Config()

	t.Run("DepositContract", func(t *testing.T) {
		_, err := n.DepositContract()
		assert.NoError(t, err)
	})
}
