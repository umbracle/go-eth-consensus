package spec

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	ssz "github.com/ferranbt/fastssz"
	"github.com/golang/snappy"
	"github.com/stretchr/testify/require"
	consensus "github.com/umbracle/go-eth-consensus"
)

func TestShuffle(t *testing.T) {
	path := []string{"shuffling", "mapping.yaml"}

	listTestData(path, func(th *testHandler) {
		var test struct {
			Seed    consensus.Root
			Count   uint64
			Mapping []uint64
		}
		th.decode(t, &test)

		for i := uint64(0); i < test.Count; i++ {
			index := ComputeShuffleIndex(i, test.Count, test.Seed)
			require.Equal(t, test.Mapping[i], index)
		}
	})
}

func listTestData(paths []string, handlerFn func(tt *testHandler)) error {
	folder := "../eth2.0-spec-tests/tests/mainnet"
	regex := ".*\\/" + strings.Join(paths, "\\/.*\\/")

	fmt.Println(regex)

	r, err := regexp.Compile(regex)
	if err != nil {
		panic(err)
	}

	return filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			return nil
		}
		if r.MatchString(path) {
			handler := &testHandler{
				path: path,
			}
			handlerFn(handler)
		}
		return nil
	})
}

type testHandler struct {
	path string
}

func (th *testHandler) decodeFile(t *testing.T, subPath string, obj interface{}) {
	path := th.path
	if subPath != "" {
		path = filepath.Join(path, subPath)
	}
	var content []byte
	ok, err := fileExists(path)
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		content, err = ioutil.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		if err := ssz.UnmarshalSSZTest(content, obj); err != nil {
			t.Fatal(err)
		}
	} else {
		// try to read the file as ssz_snappy
		snappyPath := path + ".ssz_snappy"

		ok, err := fileExists(snappyPath)
		if err != nil {
			t.Fatal(err)
		}
		if !ok {
			t.Fatalf("file %s not found (neither snappy)", subPath)
		}

		snappyContent, err := ioutil.ReadFile(snappyPath)
		if err != nil {
			t.Fatal(err)
		}
		content, err = snappy.Decode(nil, snappyContent)
		if err != nil {
			t.Fatal(err)
		}

		sszObj, ok := obj.(ssz.Unmarshaler)
		if !ok {
			t.Fatalf("obj '%s' is not ssz for snappy decompress", subPath)
		}
		if err := sszObj.UnmarshalSSZ(content); err != nil {
			t.Fatal(err)
		}
	}
}

func (th *testHandler) decode(t *testing.T, obj interface{}) {
	th.decodeFile(t, "", obj)
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
