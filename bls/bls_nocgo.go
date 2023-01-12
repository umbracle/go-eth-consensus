//go:build !cgo
// +build !cgo

package bls

import (
	"fmt"
	"math/big"

	bls12381 "github.com/kilic/bls12-381"
)

type blstPublicKey = bls12381.PointG1
type blstSignature = bls12381.PointG2

var domain = []byte("BLS_SIG_BLS12381G2_XMD:SHA-256_SSWU_RO_POP_")

// Signature is a Bls signature
type Signature struct {
	sig *blstSignature
}

func (s *Signature) Deserialize(buf []byte) error {
	g2, err := bls12381.NewG2().FromBytes(buf)
	if err != nil {
		return err
	}
	s.sig = g2
	return nil
}

func (s *Signature) Serialize() (res [96]byte) {
	buf := bls12381.NewG2().ToCompressed(s.sig)
	copy(res[:], buf)
	return
}

func (s *Signature) VerifyByte(pub *PublicKey, msg []byte) bool {
	g2, err := bls12381.NewG2().HashToCurve(msg, domain)
	if err != nil {
		panic(err)
	}
	fmt.Println(g2)

	e := bls12381.NewEngine()
	e.AddPair(pub.pub, g2)

	return e.Result().IsOne()
}

// PublicKey is a Bls public key
type PublicKey struct {
	pub *blstPublicKey
}

func (p *PublicKey) Deserialize(buf []byte) error {
	g1, err := bls12381.NewG1().FromBytes(buf)
	if err != nil {
		return err
	}
	p.pub = g1
	return nil
}

func (p *PublicKey) Serialize() (res [48]byte) {
	buf := bls12381.NewG1().ToCompressed(p.pub)
	copy(res[:], buf)
	return
}

// SecretKey is a Bls secret key
type SecretKey struct {
	key *big.Int
}

func (s *SecretKey) Unmarshal(data []byte) error {
	s.key = new(big.Int).SetBytes(data)
	return nil
}

func (s *SecretKey) Marshal() ([]byte, error) {
	return s.key.Bytes(), nil
}

func (s *SecretKey) GetPublicKey() *PublicKey {
	p := new(blstPublicKey)
	p = bls12381.NewG1().MulScalarBig(p, &bls12381.G1One, s.key)
	return &PublicKey{pub: p}
}

func (s *SecretKey) Sign(msg []byte) *Signature {
	g2, err := bls12381.NewG2().HashToCurve(msg, domain)
	if err != nil {
		panic(err)
	}

	g2 = bls12381.NewG2().MulScalarBig(g2, g2, s.key)
	return &Signature{sig: g2}
}

func RandomKey() *SecretKey {
	return nil
}
