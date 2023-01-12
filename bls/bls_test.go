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

	sig0, err := priv.Sign(msg)
	require.NoError(t, err)

	sig0B := sig0.Serialize()

	sig1 := &Signature{}
	if err := sig1.Deserialize(sig0B[:]); err != nil {
		t.Fatal(err)
	}

	pub0 := priv.GetPublicKey().Serialize()
	pub1 := &PublicKey{}
	if err := pub1.Deserialize(pub0[:]); err != nil {
		t.Fatal(err)
	}

	valid, err := sig1.VerifyByte(pub1, msg)
	require.NoError(t, err)
	require.True(t, valid)
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
		if err = pub.Deserialize(pubKeyB); err != nil {
			if !obj.Output {
				return
			}
			t.Fatal("failed to unmarshal pubkey")
		}

		sig := new(Signature)
		if err := sig.Deserialize(signatureB); err != nil {
			if !obj.Output {
				return
			}
			t.Fatal("failed to unmarshal signature")
		}

		output, err := sig.VerifyByte(pub, messageB)
		require.NoError(t, err)
		require.Equal(t, obj.Output, output)
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

		sig, err := sec.Sign(messageB)
		require.NoError(t, err)

		// we should be able to verify a signature
		verify, err := sig.VerifyByte(sec.GetPublicKey(), messageB)
		require.NoError(t, err)
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

func BenchmarkBLS_Sign(b *testing.B) {
	msg := []byte("msg")
	priv := RandomKey()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		priv.Sign(msg)
	}
}

func BenchmarkBLS_Verify(b *testing.B) {
	msg := []byte("msg")
	priv := RandomKey()
	pub := priv.GetPublicKey()

	sign, err := priv.Sign(msg)
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sign.VerifyByte(pub, msg)
	}
}
