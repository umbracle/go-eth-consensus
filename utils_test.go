package consensus

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUint256Encoding(t *testing.T) {
	v := Uint256{0x1}
	out, err := v.MarshalText()
	assert.NoError(t, err)
	assert.Equal(t, out, []byte("1"))

	vv := Uint256{}
	assert.NoError(t, vv.UnmarshalText(out))
	assert.Equal(t, v, vv)
}
