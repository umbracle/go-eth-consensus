package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuilderEndpoint(t *testing.T) {
	n := New("http://127.0.0.1:4011").Builder()

	t.Run("RegisterValidator", func(t *testing.T) {
		err := n.RegisterValidator()
		assert.NoError(t, err)
	})
}
