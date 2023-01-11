//go:build !cgo
// +build !cgo

package bls

import (
	bls12381 "github.com/kilic/bls12-381"
)

type blstPublicKey = blst.PointG1
type blstSignature = blst.PointG2

// Signature is a Bls signature
type Signature struct {
	sig *blstSignature
}

func (s *Signature) Deserialize(buf []byte) error {
	return nil
}

func (s *Signature) Serialize() (buf [96]byte) {
	return
}

func (s *Signature) VerifyByte(pub *PublicKey, msg []byte) bool {
	return false
}

// PublicKey is a Bls public key
type PublicKey struct {
	pub *blstPublicKey
}

func (p *PublicKey) Deserialize(buf []byte) error {
	return nil
}

func (p *PublicKey) Serialize() (res [48]byte) {
	return
}

// SecretKey is a Bls secret key
type SecretKey struct {
}

func (s *SecretKey) Unmarshal(data []byte) error {
	return nil
}

func (s *SecretKey) Marshal() ([]byte, error) {
	return nil, nil
}

func (s *SecretKey) GetPublicKey() *PublicKey {
	return &PublicKey{}
}

func (s *SecretKey) Sign(msg []byte) *Signature {
	return &Signature{}
}

func RandomKey() *SecretKey {
	return nil
}
