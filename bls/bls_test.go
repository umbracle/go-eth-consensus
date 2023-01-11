package bls

import (
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBLS_Simple(t *testing.T) {
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

func TestBLS_Verify(t *testing.T) {
	type ref struct {
		Input struct {
			Pubkey    string
			Message   string
			Signature string
		}
		Output bool
	}

	readBLSDir(t, "verify", new(ref), func(i interface{}) {
		obj := i.(*ref)

		pubKeyB, err := hex.DecodeString(obj.Input.Pubkey[2:])
		require.NoError(t, err)

		messageB, err := hex.DecodeString(obj.Input.Message[2:])
		require.NoError(t, err)

		signatureB, err := hex.DecodeString(obj.Input.Signature[2:])
		require.NoError(t, err)

		pub := new(PublicKey)
		require.NoError(t, pub.Deserialize(pubKeyB))

		sig := new(Signature)
		require.NoError(t, sig.Deserialize(signatureB))

		output := sig.VerifyByte(pub, messageB)
		require.Equal(t, output, obj.Output)
	})
}

func TestBLS_Sign(t *testing.T) {
	type ref struct {
		Input struct {
			Privkey string
			Message string
		}
		Output *string
	}

	readBLSDir(t, "sign", new(ref), func(i interface{}) {
		obj := i.(*ref)

		privKeyB, err := hex.DecodeString(obj.Input.Privkey[2:])
		require.NoError(t, err)

		messageB, err := hex.DecodeString(obj.Input.Message[2:])
		require.NoError(t, err)

		if obj.Output == nil {
			return
		}

		sec := new(SecretKey)
		require.NoError(t, sec.Unmarshal(privKeyB))

		sig := sec.Sign(messageB)

		// we should be able to verify a signature
		verify := sig.VerifyByte(sec.GetPublicKey(), messageB)
		require.True(t, verify)

		outputB, err := hex.DecodeString((*obj.Output)[2:])
		require.NoError(t, err)

		sigB := sig.Serialize()
		require.Equal(t, outputB, sigB[:])
	})
}

func readBLSDir(t *testing.T, path string, ref interface{}, callback func(interface{})) {
	fullPath := filepath.Join("../eth2.0-spec-tests/bls", path)

	files, err := ioutil.ReadDir(fullPath)
	require.NoError(t, err)

	for _, file := range files {
		data, err := ioutil.ReadFile(filepath.Join(fullPath, file.Name()))
		require.NoError(t, err)

		obj := reflect.New(reflect.TypeOf(ref).Elem()).Interface()

		err = json.Unmarshal(data, obj)
		require.NoError(t, err)

		callback(obj)
	}
}
