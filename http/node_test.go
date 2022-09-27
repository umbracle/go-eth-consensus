package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNodeEndpoint(t *testing.T) {
	n := New("http://127.0.0.1:4010", WithUntrackedKeys()).Node()

	t.Run("Identity", func(t *testing.T) {
		_, err := n.Identity()
		assert.NoError(t, err)
	})

	t.Run("Peers", func(t *testing.T) {
		_, err := n.Peers()
		assert.NoError(t, err)
	})

	t.Run("GetPeer", func(t *testing.T) {
		_, err := n.GetPeer("a")
		assert.NoError(t, err)
	})

	t.Run("PeerCount", func(t *testing.T) {
		_, err := n.PeerCount()
		assert.NoError(t, err)
	})

	t.Run("Version", func(t *testing.T) {
		_, err := n.Version()
		assert.NoError(t, err)
	})

	t.Run("Syncing", func(t *testing.T) {
		_, err := n.Syncing()
		assert.NoError(t, err)
	})

	t.Run("Health", func(t *testing.T) {
		status, err := n.Health()
		assert.NoError(t, err)
		assert.True(t, status)
	})
}
