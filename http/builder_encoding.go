// Code generated by fastssz. DO NOT EDIT.
// Hash: fba4ba076684515d3fe1274c2c488936a0aa8ed5eb824fb848bad138b07616a9
package http

import (
	ssz "github.com/ferranbt/fastssz"
)

// MarshalSSZ ssz marshals the RegisterValidatorRequest object
func (r *RegisterValidatorRequest) MarshalSSZ() ([]byte, error) {
	return ssz.MarshalSSZ(r)
}

// MarshalSSZTo ssz marshals the RegisterValidatorRequest object to a target array
func (r *RegisterValidatorRequest) MarshalSSZTo(buf []byte) (dst []byte, err error) {
	dst = buf

	// Field (0) 'FeeRecipient'
	dst = append(dst, r.FeeRecipient[:]...)

	// Field (1) 'GasLimit'
	dst = ssz.MarshalUint64(dst, r.GasLimit)

	// Field (2) 'Timestamp'
	dst = ssz.MarshalUint64(dst, r.Timestamp)

	// Field (3) 'Pubkey'
	dst = append(dst, r.Pubkey[:]...)

	return
}

// UnmarshalSSZ ssz unmarshals the RegisterValidatorRequest object
func (r *RegisterValidatorRequest) UnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size != 84 {
		return ssz.ErrSize
	}

	// Field (0) 'FeeRecipient'
	copy(r.FeeRecipient[:], buf[0:20])

	// Field (1) 'GasLimit'
	r.GasLimit = ssz.UnmarshallUint64(buf[20:28])

	// Field (2) 'Timestamp'
	r.Timestamp = ssz.UnmarshallUint64(buf[28:36])

	// Field (3) 'Pubkey'
	copy(r.Pubkey[:], buf[36:84])

	return err
}

// SizeSSZ returns the ssz encoded size in bytes for the RegisterValidatorRequest object
func (r *RegisterValidatorRequest) SizeSSZ() (size int) {
	size = 84
	return
}

// HashTreeRoot ssz hashes the RegisterValidatorRequest object
func (r *RegisterValidatorRequest) HashTreeRoot() ([32]byte, error) {
	return ssz.HashWithDefaultHasher(r)
}

// HashTreeRootWith ssz hashes the RegisterValidatorRequest object with a hasher
func (r *RegisterValidatorRequest) HashTreeRootWith(hh ssz.HashWalker) (err error) {
	indx := hh.Index()

	// Field (0) 'FeeRecipient'
	hh.PutBytes(r.FeeRecipient[:])

	// Field (1) 'GasLimit'
	hh.PutUint64(r.GasLimit)

	// Field (2) 'Timestamp'
	hh.PutUint64(r.Timestamp)

	// Field (3) 'Pubkey'
	hh.PutBytes(r.Pubkey[:])

	hh.Merkleize(indx)
	return
}

// GetTree ssz hashes the RegisterValidatorRequest object
func (r *RegisterValidatorRequest) GetTree() (*ssz.Node, error) {
	return ssz.ProofTree(r)
}
