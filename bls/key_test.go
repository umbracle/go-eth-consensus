package bls

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKey_Encoding(t *testing.T) {
	key := NewRandomKey()

	data, err := key.Marshal()
	assert.NoError(t, err)

	key1 := &Key{}
	assert.NoError(t, key1.Unmarshal(data))

	assert.True(t, key.Equal(key1))
}

func TestKey_Keystore_Fixture(t *testing.T) {
	content, err := ioutil.ReadFile("./fixtures/keystore.json")
	assert.NoError(t, err)

	var data []json.RawMessage
	assert.NoError(t, json.Unmarshal(content, &data))

	_, err = FromKeystore(data[0], "𝔱𝔢𝔰𝔱𝔭𝔞𝔰𝔰𝔴𝔬𝔯𝔡🔑")
	assert.NoError(t, err)
}
