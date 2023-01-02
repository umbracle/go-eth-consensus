package spec

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"

	consensus "github.com/umbracle/go-eth-consensus"
)

const shuffleRoundCount = 90

func ComputeShuffleIndex(index, indexCount uint64, seed consensus.Root) uint64 {
	if index >= indexCount {
		panic(fmt.Sprintf("BAD: index %d higher than count %d", index, indexCount))
	}

	for i := 0; i < shuffleRoundCount; i++ {
		input := append(seed[:], byte(i))
		hash := sha256.New()
		hash.Write(input)

		hashValue := binary.LittleEndian.Uint64(hash.Sum(nil)[:8])

		pivot := hashValue % indexCount
		flip := (pivot + indexCount - index) % indexCount

		position := index
		if flip > index {
			position = flip
		}

		positionByteArray := make([]byte, 4)
		binary.LittleEndian.PutUint32(positionByteArray, uint32(position>>8))
		input2 := append(seed[:], byte(i))
		input2 = append(input2, positionByteArray...)

		hash.Reset()
		hash.Write(input2)

		source := hash.Sum(nil)
		byteVal := source[(position%256)/8]
		bitVal := (byteVal >> (position % 8)) % 2
		if bitVal == 1 {
			index = flip
		}
	}

	return index
}
