package spec

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"

	consensus "github.com/umbracle/go-eth-consensus"
)

const shuffleRoundCount = 90

func computeShuffleIndex(index, indexCount uint64, seed consensus.Root) uint64 {
	if index >= indexCount {
		panic(fmt.Sprintf("BAD: index %d higher than count %d", index, indexCount))
	}

	for i := 0; i < shuffleRoundCount; i++ {
		input := make([]byte, 0, len(seed)+1)
		input = append(input, seed[:]...)
		input = append(input, byte(i))

		res := sha256.Sum256(input)
		hashValue := binary.LittleEndian.Uint64(res[:8])
		pivot := hashValue % indexCount
		flip := (pivot + indexCount - index) % indexCount

		position := index
		if flip > index {
			position = flip
		}

		positionByteArray := make([]byte, 4)
		binary.LittleEndian.PutUint32(positionByteArray, uint32(position>>8))

		input2 := make([]byte, 0, len(seed)+5)
		input2 = append(input2, seed[:]...)
		input2 = append(input2, byte(i))
		input2 = append(input2, positionByteArray...)

		source := sha256.Sum256(input2)
		byteVal := source[(position%256)/8]
		bitVal := (byteVal >> (position % 8)) % 2
		if bitVal == 1 {
			index = flip
		}
	}

	return index
}
