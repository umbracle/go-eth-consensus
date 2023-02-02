//go:build cgo
// +build cgo

package bls

import (
	"crypto/rand"
	"fmt"

	blst "github.com/supranational/blst/bindings/go"
)

type blstPublicKey = blst.P1Affine
type blstSignature = blst.P2Affine

var dst = []byte("BLS_SIG_BLS12381G2_XMD:SHA-256_SSWU_RO_POP_")

// Signature is a Bls signature
type Signature struct {
	sig *blstSignature
}

func (s *Signature) Deserialize(buf []byte) error {
	sig := new(blstSignature).Uncompress(buf)
	if sig == nil {
		return fmt.Errorf("failed to deserialize")
	}
	if !sig.SigValidate(false) {
		return fmt.Errorf("signature not in group")
	}
	s.sig = sig
	return nil
}

func (s *Signature) Serialize() (buf [96]byte) {
	copy(buf[:], s.sig.Compress())
	return
}

func (s *Signature) VerifyByte(pub *PublicKey, msg []byte) (bool, error) {
	return s.sig.Verify(false, pub.pub, false, msg, dst), nil
}

func (s *Signature) FastAggregateVerify(pubKeys []*PublicKey, msg []byte) (bool, error) {
	raw := make([]*blstPublicKey, len(pubKeys))
	for indx, i := range pubKeys {
		raw[indx] = i.pub
	}

	return s.sig.FastAggregateVerify(true, raw, msg, dst), nil
}

func AggregateSignatures(sigs []*Signature) *Signature {
	if len(sigs) == 0 {
		return nil
	}
	raw := make([]*blstSignature, len(sigs))
	for indx, i := range sigs {
		raw[indx] = i.sig
	}

	sig := new(blst.P2Aggregate)
	sig.Aggregate(raw, false)

	return &Signature{sig: sig.ToAffine()}
}

// PublicKey is a Bls public key
type PublicKey struct {
	pub *blstPublicKey
}

func (p *PublicKey) Deserialize(buf []byte) error {
	pub := new(blstPublicKey).Uncompress(buf)
	if pub == nil {
		return fmt.Errorf("failed to deserialize")
	}
	if !pub.KeyValidate() {
		return fmt.Errorf("point at infinity")
	}
	p.pub = pub
	return nil
}

func (p *PublicKey) Serialize() (res [48]byte) {
	copy(res[:], p.pub.Compress())
	return
}

// SecretKey is a Bls secret key
type SecretKey struct {
	key *blst.SecretKey
}

func (s *SecretKey) Unmarshal(data []byte) error {
	s.key = new(blst.SecretKey).Deserialize(data)
	return nil
}

func (s *SecretKey) Marshal() ([]byte, error) {
	return s.key.Serialize(), nil
}

func (s *SecretKey) GetPublicKey() *PublicKey {
	pub := new(blstPublicKey).From(s.key)
	return &PublicKey{pub: pub}
}

func (s *SecretKey) Sign(msg []byte) (*Signature, error) {
	sig := new(blstSignature).Sign(s.key, msg, dst)
	return &Signature{sig: sig}, nil
}

func RandomKey() *SecretKey {
	var ikm [32]byte
	_, _ = rand.Read(ikm[:])
	sk := blst.KeyGen(ikm[:])

	sec := &SecretKey{
		key: sk,
	}
	return sec
}
