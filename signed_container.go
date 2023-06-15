package consensus

import ssz "github.com/ferranbt/fastssz"

type container interface {
	ssz.Marshaler
	ssz.Unmarshaler
	ssz.HashRoot
}

type SignedContainer[T container] struct {
	Container T         `json:"message"`
	Signature Signature `json:"signature"`
}

// HashTreeRoot ssz hashes the SignedContainer object
func (s *SignedContainer[T]) HashTreeRoot() ([32]byte, error) {
	return ssz.HashWithDefaultHasher(s)
}

// GetTree ssz hashes the SignedContainer object
func (s *SignedContainer[T]) GetTree() (*ssz.Node, error) {
	return ssz.ProofTree(s)
}

// SizeSSZ returns the ssz encoded size in bytes for the SignedContainer object
func (s *SignedContainer[T]) SizeSSZ() (size int) {
	size = s.Container.SizeSSZ() + 12

	return
}

// MarshalSSZ ssz marshals the SignedContainer object
func (s *SignedContainer[T]) MarshalSSZ() ([]byte, error) {
	return ssz.MarshalSSZ(s)
}

// MarshalSSZTo ssz marshals the SignedContainer object to a target array
func (s *SignedContainer[T]) MarshalSSZTo(buf []byte) (dst []byte, err error) {
	dst = buf

	// Field (0) 'Container'
	if dst, err = s.Container.MarshalSSZTo(dst); err != nil {
		return
	}

	// Field (1) 'Signature'
	dst = append(dst, s.Signature[:]...)

	return
}

func (s *SignedContainer[T]) HashTreeRootWith(hh ssz.HashWalker) (err error) {
	indx := hh.Index()

	// Field (0) 'Container'
	if err = s.Container.HashTreeRootWith(hh); err != nil {
		return
	}

	// Field (1) 'Signature'
	hh.PutBytes(s.Signature[:])

	hh.Merkleize(indx)
	return
}

// UnmarshalSSZ ssz unmarshals the SignedContainer object
func (s *SignedContainer[T]) UnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size != 112 {
		return ssz.ErrSize
	}

	// Field (0) 'Container'
	if err = s.Container.UnmarshalSSZ(buf[0:16]); err != nil {
		return err
	}

	// Field (1) 'Signature'
	copy(s.Signature[:], buf[16:112])

	return err
}
