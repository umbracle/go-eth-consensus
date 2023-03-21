package bitlist

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBitmap_SetOutOfBounds(t *testing.T) {
	b := NewBitlist(10)

	bb := b.Copy()
	b.SetBitAt(10, true)

	require.True(t, bb.Equal(b))
}

func TestBitmap_SetIndx(t *testing.T) {
	size := uint64(10)
	b := NewBitlist(uint64(size))

	// all empty
	for i := uint64(0); i < size; i++ {
		require.False(t, b.BitAt(i))
	}

	// set indexes to true
	for i := uint64(0); i < size; i++ {
		b.SetBitAt(i, true)
		require.True(t, b.BitAt(i))

		for j := uint64(0); j < size; j++ {
			require.Equal(t, b.BitAt(j), j <= i)
		}
	}

	// all indexes are full
	for i := uint64(0); i < size; i++ {
		require.True(t, b.BitAt(i))
	}

	// set indexes to false
	for i := uint64(0); i < size; i++ {
		b.SetBitAt(i, false)
		require.False(t, b.BitAt(i))

		for j := uint64(0); j < size; j++ {
			require.Equal(t, b.BitAt(j), j > i)
		}
	}
}
