package bls

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"

	uuid "github.com/hashicorp/go-uuid"
	"github.com/umbracle/ethgo/keystore"
)

// Key is a reference to a key in the keymanager
type Key struct {
	Id  string
	Pub *PublicKey
	Prv *SecretKey
}

func (k *Key) Unmarshal(data []byte) error {
	k.Prv = &SecretKey{}
	if err := k.Prv.Unmarshal(data); err != nil {
		return err
	}
	k.Pub = k.Prv.GetPublicKey()
	return nil
}

func (k *Key) Marshal() ([]byte, error) {
	return k.Prv.Marshal()
}

func (k *Key) Equal(kk *Key) bool {
	a := k.Pub.Serialize()
	b := kk.Pub.Serialize()
	return bytes.Equal(a[:], b[:])
}

func (k *Key) PubKey() (out [48]byte) {
	return k.Pub.Serialize()
}

func (k *Key) Sign(root [32]byte) ([96]byte, error) {
	signed, err := k.Prv.Sign(root[:])
	if err != nil {
		return [96]byte{}, err
	}
	return signed.Serialize(), nil
}

func NewKeyFromPriv(priv []byte) (*Key, error) {
	k := &Key{}
	if err := k.Unmarshal(priv); err != nil {
		return nil, err
	}
	return k, nil
}

func NewRandomKey() *Key {
	sec := RandomKey()
	id, _ := uuid.GenerateUUID()

	k := &Key{
		Id:  id,
		Prv: sec,
		Pub: sec.GetPublicKey(),
	}
	return k
}

func FromKeystore(content []byte, password string) (*Key, error) {
	var dec map[string]interface{}
	if err := json.Unmarshal(content, &dec); err != nil {
		return nil, err
	}

	priv, err := keystore.DecryptV4(content, password)
	if err != nil {
		return nil, err
	}
	key, err := NewKeyFromPriv(priv)
	if err != nil {
		return nil, err
	}

	pub := key.PubKey()
	if hex.EncodeToString(pub[:]) != dec["pubkey"] {
		return nil, fmt.Errorf("pub key does not match")
	}
	key.Id = dec["uuid"].(string)
	return key, nil
}

func ToKeystore(k *Key, password string) ([]byte, error) {
	priv, err := k.Prv.Marshal()
	if err != nil {
		return nil, err
	}

	keystore, err := keystore.EncryptV4(priv, password)
	if err != nil {
		return nil, err
	}

	var dec map[string]interface{}
	if err := json.Unmarshal(keystore, &dec); err != nil {
		return nil, err
	}

	serializePub := k.Pub.Serialize()

	dec["pubkey"] = hex.EncodeToString(serializePub[:])
	dec["uuid"] = k.Id

	// small error here, params is set to nil in ethgo
	a := dec["crypto"].(map[string]interface{})
	b := a["checksum"].(map[string]interface{})
	b["params"] = map[string]interface{}{}

	return json.Marshal(dec)
}
