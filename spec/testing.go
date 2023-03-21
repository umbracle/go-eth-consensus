package spec

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	ssz "github.com/ferranbt/fastssz"
	"github.com/golang/snappy"
	"github.com/stretchr/testify/require"
)

var (
	testDataFolder = "../eth2.0-spec-tests/tests"
)

func listTestData(t *testing.T, path string, handlerFn func(tt *testHandler)) {
	matches, err := filepath.Glob(filepath.Join(testDataFolder, path))
	require.NoError(t, err)

	if len(matches) == 0 {
		t.Fatal("no matches found")
	}

	for _, m := range matches {
		//t.Run(m, func(t *testing.T) {
		handler := &testHandler{
			t:    t,
			path: m,
		}
		handlerFn(handler)
		//})
	}
}

type testHandler struct {
	t    *testing.T
	path string
}

func (th *testHandler) decodeFile(subPath string, obj interface{}, maybeEmpty ...bool) bool {
	path := filepath.Join(th.path, subPath)
	var content []byte

	ok, err := fileExists(path)
	require.NoError(th.t, err)

	if ok {
		content, err = ioutil.ReadFile(path)
		require.NoError(th.t, err)

		err = ssz.UnmarshalSSZTest(content, obj)
		require.NoError(th.t, err)
	} else {
		// try to read the file as ssz_snappy
		snappyPath := path + ".ssz_snappy"

		ok, err := fileExists(snappyPath)
		require.NoError(th.t, err)

		if !ok {
			if len(maybeEmpty) != 0 && maybeEmpty[0] {
				// the file might not exist
				return false
			}
			th.t.Fatalf("file '%s' not found (neither snappy)", subPath)
		}

		snappyContent, err := ioutil.ReadFile(snappyPath)
		require.NoError(th.t, err)

		content, err = snappy.Decode(nil, snappyContent)
		require.NoError(th.t, err)

		sszObj, ok := obj.(ssz.Unmarshaler)
		if !ok {
			th.t.Fatalf("obj '%s' is not ssz for snappy decompress", subPath)
		}
		err = sszObj.UnmarshalSSZ(content)
		require.NoError(th.t, err)
	}

	return true
}

func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
