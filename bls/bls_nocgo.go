//go:build !cgo
// +build !cgo

package bls

import (
	"crypto/rand"
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
	g2, err := bls12381.NewG2().FromCompressed(buf)
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

func (s *Signature) VerifyByte(pub *PublicKey, msg []byte) (bool, error) {
	return s.verifyImpl(pub.pub, msg)
}

func (s *Signature) verifyImpl(g1 *bls12381.PointG1, msg []byte) (bool, error) {
	hash, err := bls12381.NewG2().HashToCurve(msg, domain)
	if err != nil {
		return false, err
	}

	e := bls12381.NewEngine()
	e.AddPairInv(e.G1.One(), s.sig)
	e.AddPair(g1, hash)

	return e.Check(), nil
}

func (s *Signature) FastAggregateVerify(pubKeys []*PublicKey, msg []byte) (bool, error) {
	if bls12381.NewG2().IsZero(s.sig) {
		// signature is infinite
		return false, nil
	}

	// aggregate public keys
	aggPub := new(bls12381.PointG1)
	g1 := bls12381.NewG1()

	for _, pub := range pubKeys {
		aggPub = g1.Add(aggPub, aggPub, pub.pub)
	}

	ok, err := s.verifyImpl(aggPub, msg)
	if err != nil {
		return false, err
	}
	return ok, nil
}

func AggregateSignatures(sigs []*Signature) *Signature {
	if len(sigs) == 0 {
		return nil
	}

	aggSig := new(bls12381.PointG2)
	g2 := bls12381.NewG2()

	for _, sig := range sigs {
		aggSig = g2.Add(aggSig, aggSig, sig.sig)
	}

	return &Signature{sig: aggSig}
}

// PublicKey is a Bls public key
type PublicKey struct {
	pub *blstPublicKey
}

func (p *PublicKey) Deserialize(buf []byte) error {
	g1, err := bls12381.NewG1().FromCompressed(buf)
	if err != nil {
		return err
	}
	if bls12381.NewG1().IsZero(g1) {
		return fmt.Errorf("infinity")
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

func (s *SecretKey) Sign(msg []byte) (*Signature, error) {
	hash, err := bls12381.NewG2().HashToCurve(msg, domain)
	if err != nil {
		return nil, err
	}

	g2 := bls12381.NewG2()
	g2.MulScalarBig(hash, hash, s.key)

	return &Signature{sig: hash}, nil
}

var curveOrder, _ = new(big.Int).SetString("73eda753299d7d483339d80809a1d80553bda402fffe5bfeffffffff00000001", 16)

func RandomKey() *SecretKey {
	k, err := rand.Int(rand.Reader, curveOrder)
	if err != nil {
		panic(err)
	}
	return &SecretKey{key: k}
}
