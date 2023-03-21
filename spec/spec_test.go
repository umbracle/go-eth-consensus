package spec

import (
	_ "embed"
	"fmt"
	"strings"
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

//go:embed presets/phase0.yaml
var mainnetPresetPhase0 []byte

func TestPresetMainnet(t *testing.T) {
	var out map[string]interface{}
	require.NoError(t, yaml.Unmarshal(mainnetPresetPhase0, &out))

	var specOut map[string]interface{}
	require.NoError(t, mapstructure.Decode(Spec, &specOut))

	// check that the preset values from 'out' match
	// the value from the spec struct
	for yamlKey, presetVal := range out {
		// convert the key name (i.e. ABC_DEF) to cammel case (i.e. AbcDef)
		lowerKey := strings.ToLower(yamlKey)
		var key string

		for i := 0; i < len(yamlKey); i++ {
			if string(lowerKey[i]) == "_" {
				key += strings.ToTitle(string(lowerKey[i+1]))
				i++
			} else {
				key += string(lowerKey[i])
			}
		}

		key = strings.Title(key)

		if val, ok := specOut[key]; ok {
			// avoid uint8 and int comparisions by using the fmt.Sprinot to format
			require.Equal(t, fmt.Sprintf("%d", val), fmt.Sprintf("%d", presetVal))
		}
	}
}
