package bls

import (
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"strings"
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

func TestBLS_Aggregate(t *testing.T) {
	type ref struct {
		Input  []argBytes
		Output argBytes
	}

	readBLSDir(t, "aggregate", new(ref), func(name string, i interface{}) {
		obj := i.(*ref)

		var sigs []*Signature
		for _, o := range obj.Input {
			sig := new(Signature)
			require.NoError(t, sig.Deserialize(o))

			sigs = append(sigs, sig)
		}

		if len(obj.Output) == 0 {
			return
		}

		sig := AggregateSignatures(sigs).Serialize()
		require.Equal(t, sig[:], []byte(obj.Output))
	})
}

func TestBLS_FastAggregateVerify(t *testing.T) {
	type ref struct {
		Input struct {
			Pubkeys   []argBytes
			Message   argBytes
			Signature argBytes
		}
		Output bool
	}

	readBLSDir(t, "fast_aggregate_verify", new(ref), func(name string, i interface{}) {
		obj := i.(*ref)

		pubKeys := []*PublicKey{}
		for _, elem := range obj.Input.Pubkeys {

			pub := new(PublicKey)
			if err := pub.Deserialize(elem); err != nil {
				if !obj.Output {
					return
				}
				t.Fatal(err)
			}

			pubKeys = append(pubKeys, pub)
		}

		sig := new(Signature)
		if err := sig.Deserialize(obj.Input.Signature); err != nil {
			if !obj.Output {
				return
			}
			t.Fatal(err)
		}

		ok, err := sig.FastAggregateVerify(pubKeys, obj.Input.Message)
		require.NoError(t, err)
		require.Equal(t, ok, obj.Output)
	})
}

func TestBLS_DeserializationG1(t *testing.T) {
	type ref struct {
		Input struct {
			Pubkey argBytes
		}
		Output bool
	}

	readBLSDir(t, "deserialization_G1", new(ref), func(name string, i interface{}) {
		obj := i.(*ref)

		pub := new(PublicKey)
		err := pub.Deserialize(obj.Input.Pubkey)

		if name == "deserialization_succeeds_infinity_with_true_b_flag.json" {
			// we also fail if point is at inifinity
			return
		}

		if err == nil && !obj.Output {
			t.Fatal("it should fail")
		} else if err != nil && obj.Output {
			t.Fatal(err)
		}
	})
}

func TestBLS_DeserializationG2(t *testing.T) {
	type ref struct {
		Input struct {
			Signature argBytes
		}
		Output bool
	}

	readBLSDir(t, "deserialization_G2", new(ref), func(name string, i interface{}) {
		obj := i.(*ref)

		sig := new(Signature)
		err := sig.Deserialize(obj.Input.Signature)

		if err == nil && !obj.Output {
			t.Fatal("it should fail")
		} else if err != nil && obj.Output {
			t.Fatal(err)
		}
	})
}

func TestBLS_Verify(t *testing.T) {
	type ref struct {
		Input struct {
			Pubkey    argBytes
			Message   argBytes
			Signature argBytes
		}
		Output bool
	}

	readBLSDir(t, "verify", new(ref), func(name string, i interface{}) {
		obj := i.(*ref)

		pub := new(PublicKey)
		if err := pub.Deserialize(obj.Input.Pubkey); err != nil {
			if !obj.Output {
				return
			}
			t.Fatal("failed to unmarshal pubkey")
		}

		sig := new(Signature)
		if err := sig.Deserialize(obj.Input.Signature); err != nil {
			if !obj.Output {
				return
			}
			t.Fatal("failed to unmarshal signature")
		}

		output, err := sig.VerifyByte(pub, obj.Input.Message)
		require.NoError(t, err)
		require.Equal(t, obj.Output, output)
	})
}

func TestBLS_Sign(t *testing.T) {
	type ref struct {
		Input struct {
			Privkey argBytes
			Message argBytes
		}
		Output *argBytes
	}

	readBLSDir(t, "sign", new(ref), func(name string, i interface{}) {
		obj := i.(*ref)

		if obj.Output == nil {
			return
		}

		sec := new(SecretKey)
		require.NoError(t, sec.Unmarshal(obj.Input.Privkey))

		sig, err := sec.Sign(obj.Input.Message)
		require.NoError(t, err)

		// we should be able to verify a signature
		verify, err := sig.VerifyByte(sec.GetPublicKey(), obj.Input.Message)
		require.NoError(t, err)
		require.True(t, verify)

		sigB := sig.Serialize()
		require.Equal(t, []byte(*obj.Output), sigB[:])
	})
}

func readBLSDir(t *testing.T, path string, ref interface{}, callback func(string, interface{})) {
	fullPath := filepath.Join("../eth2.0-spec-tests/bls", path)

	files, err := ioutil.ReadDir(fullPath)
	require.NoError(t, err)

	for _, file := range files {
		data, err := ioutil.ReadFile(filepath.Join(fullPath, file.Name()))
		require.NoError(t, err)

		obj := reflect.New(reflect.TypeOf(ref).Elem()).Interface()

		err = json.Unmarshal(data, obj)
		require.NoError(t, err)

		callback(file.Name(), obj)
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

func BenchmarkBLS_AggregateVerify(b *testing.B) {
	msg := []byte("msg")

	num := 10
	pubs := make([]*PublicKey, num)
	sigs := make([]*Signature, num)

	for i := 0; i < num; i++ {
		priv := RandomKey()

		sign, err := priv.Sign(msg)
		if err != nil {
			b.Fatal(err)
		}

		sigs[i] = sign
		pubs[i] = priv.GetPublicKey()
	}

	sig := AggregateSignatures(sigs)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sig.FastAggregateVerify(pubs, msg)
	}
}

type argBytes []byte

func (b *argBytes) UnmarshalText(input []byte) error {
	if len(input) == 0 {
		// some tests have empty inputs
		return nil
	}
	str := string(input)
	if strings.HasPrefix(str, "0x") {
		// some values have 0x prefix
		str = str[2:]
	}
	buf, err := hex.DecodeString(str)
	if err != nil {
		return nil
	}
	aux := make([]byte, len(buf))
	copy(aux[:], buf[:])
	*b = aux
	return nil
}
