package consensus

import (
	"fmt"
	"testing"
)

func TestSignedContainer(t *testing.T) {
	b := boxed[*AttestationData, AttestationData]{}
	fmt.Println(b)
}

type I interface {
	HashTreeRoot() ([32]byte, error)
}

type boxed[T I, PR *I] struct {
}
