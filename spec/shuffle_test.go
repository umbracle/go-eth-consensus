package spec

import (
	"testing"

	"github.com/stretchr/testify/require"
	consensus "github.com/umbracle/go-eth-consensus"
)

func TestShuffle(t *testing.T) {
	listTestData(t, "mainnet/phase0/shuffling/*/*/*", func(th *testHandler) {
		shuffleTest := &shuffleTest{}
		shuffleTest.Decode(th)

		for i := uint64(0); i < shuffleTest.Count; i++ {
			index := ComputeShuffleIndex(i, shuffleTest.Count, shuffleTest.Seed)
			require.Equal(t, shuffleTest.Mapping[i], index)
		}
	})
}

type shuffleTest struct {
	Seed    consensus.Root
	Count   uint64
	Mapping []uint64
}

func (s *shuffleTest) Decode(th *testHandler) {
	th.decodeFile("mapping.yaml", &s)
}
