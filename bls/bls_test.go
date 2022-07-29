package bls

import (
	"testing"
)

func TestBLS(t *testing.T) {
	msg := []byte("msg")
	priv := RandomKey()

	sig0 := priv.Sign(msg).Serialize()
	sig1 := &Signature{}
	if err := sig1.Deserialize(sig0[:]); err != nil {
		t.Fatal(err)
	}

	pub0 := priv.GetPublicKey().Serialize()
	pub1 := &PublicKey{}
	if err := pub1.Deserialize(pub0[:]); err != nil {
		t.Fatal(err)
	}

	if !sig1.VerifyByte(pub1, msg) {
		t.Fatal("failed to validate deposit")
	}
}
