package consensus

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	ssz "github.com/ferranbt/fastssz"
	"github.com/golang/snappy"

	"gopkg.in/yaml.v2"
)

type codec interface {
	ssz.Marshaler
	ssz.Unmarshaler
	ssz.HashRoot
}

type codecTree interface {
	GetTreeWithWrapper(w *ssz.Wrapper) (err error)
	GetTree() (*ssz.Node, error)
}

type fork string

const (
	phase0Fork    = "phase0"
	altairFork    = "altair"
	bellatrixFork = "bellatrix"
	capellaFork   = "capella"
)

type testCallback func(f fork) codec

var codecs = map[string]testCallback{
	"AttestationData":             func(f fork) codec { return new(AttestationData) },
	"Checkpoint":                  func(f fork) codec { return new(Checkpoint) },
	"AggregateAndProof":           func(f fork) codec { return new(AggregateAndProof) },
	"Attestation":                 func(f fork) codec { return new(Attestation) },
	"AttesterSlashing":            func(f fork) codec { return new(AttesterSlashing) },
	"LightClientBootstrap":        func(f fork) codec { return new(LightClientBootstrap) },
	"LightClientFinalityUpdate":   func(f fork) codec { return new(LightClientFinalityUpdate) },
	"LightClientOptimisticUpdate": func(f fork) codec { return new(LightClientOptimisticUpdate) },
	"LightClientUpdate":           func(f fork) codec { return new(LightClientUpdate) },
	"BeaconBlock": func(f fork) codec {
		if f == capellaFork {
			return new(BeaconBlockCapella)
		} else if f == altairFork {
			return new(BeaconBlockAltair)
		} else if f == bellatrixFork {
			return new(BeaconBlockBellatrix)
		}
		return new(BeaconBlockPhase0)
	},
	"BeaconBlockBody": func(f fork) codec {
		if f == capellaFork {
			return new(BeaconBlockBodyCapella)
		} else if f == altairFork {
			return new(BeaconBlockBodyAltair)
		} else if f == bellatrixFork {
			return new(BeaconBlockBodyBellatrix)
		}
		return new(BeaconBlockBodyPhase0)
	},
	"BeaconBlockHeader":  func(f fork) codec { return new(BeaconBlockHeader) },
	"Deposit":            func(f fork) codec { return new(Deposit) },
	"DepositData":        func(f fork) codec { return new(DepositData) },
	"DepositMessage":     func(f fork) codec { return new(DepositMessage) },
	"Eth1Data":           func(f fork) codec { return new(Eth1Data) },
	"Fork":               func(f fork) codec { return new(Fork) },
	"IndexedAttestation": func(f fork) codec { return new(IndexedAttestation) },
	"PendingAttestation": func(f fork) codec { return new(PendingAttestation) },
	"ProposerSlashing":   func(f fork) codec { return new(ProposerSlashing) },
	"SignedBeaconBlock": func(f fork) codec {
		if f == capellaFork {
			return new(SignedBeaconBlockCapella)
		} else if f == altairFork {
			return new(SignedBeaconBlockAltair)
		} else if f == bellatrixFork {
			return new(SignedBeaconBlockBellatrix)
		}
		return new(SignedBeaconBlockPhase0)
	},
	"SignedBeaconBlockHeader":     func(f fork) codec { return new(SignedBeaconBlockHeader) },
	"SignedVoluntaryExit":         func(f fork) codec { return new(SignedVoluntaryExit) },
	"SigningRoot":                 func(f fork) codec { return new(SigningRoot) },
	"Validator":                   func(f fork) codec { return new(Validator) },
	"VoluntaryExit":               func(f fork) codec { return new(VoluntaryExit) },
	"SyncCommittee":               func(f fork) codec { return new(SyncCommittee) },
	"SyncAggregate":               func(f fork) codec { return new(SyncAggregate) },
	"SyncCommitteeMessage":        func(f fork) codec { return new(SyncCommitteeMessage) },
	"SyncCommitteeContribution":   func(f fork) codec { return new(SyncCommitteeContribution) },
	"SignedContributionAndProof":  func(f fork) codec { return new(SignedContributionAndProof) },
	"ContributionAndProof":        func(f fork) codec { return new(ContributionAndProof) },
	"Eth1Block":                   func(f fork) codec { return new(Eth1Block) },
	"SyncAggregatorSelectionData": func(f fork) codec { return new(SyncAggregatorSelectionData) },
	"SigningData":                 func(f fork) codec { return new(SigningData) },
	"ForkData":                    func(f fork) codec { return new(ForkData) },
	"SignedAggregateAndProof":     func(f fork) codec { return new(SignedAggregateAndProof) },
	"PowBlock":                    func(f fork) codec { return new(PowBlock) },
	"ExecutionPayload": func(f fork) codec {
		if f == capellaFork {
			return new(ExecutionPayloadCapella)
		}
		return new(ExecutionPayload)
	},
	"ExecutionPayloadHeader": func(f fork) codec {
		if f == capellaFork {
			return new(ExecutionPayloadHeaderCapella)
		}
		return new(ExecutionPayloadHeader)
	},
	"BeaconState": func(f fork) codec {
		if f == altairFork {
			return new(BeaconStateAltair)
		} else if f == bellatrixFork {
			return new(BeaconStateBellatrix)
		} else if f == capellaFork {
			return new(BeaconStateCapella)
		}
		return new(BeaconStatePhase0)
	},
	"BLSToExecutionChange":       func(f fork) codec { return new(BLSToExecutionChange) },
	"HistoricalSummary":          func(f fork) codec { return new(HistoricalSummary) },
	"SignedBLSToExecutionChange": func(f fork) codec { return new(SignedBLSToExecutionChange) },
	"Withdrawal":                 func(f fork) codec { return new(Withdrawal) },
}

func testFork(t *testing.T, fork fork) {
	files := readDir(t, filepath.Join(testsPath, "/mainnet/"+string(fork)+"/ssz_static"))
	for _, f := range files {
		spl := strings.Split(f, "/")
		name := spl[len(spl)-1]

		base, ok := codecs[name]
		if !ok {
			t.Logf("type %s not found in fork %s", name, fork)
			continue
		}

		t.Run(name, func(t *testing.T) {
			files := readDir(t, filepath.Join(f, "ssz_random"))
			for _, f := range files {
				checkSSZEncoding(t, fork, f, name, base)
			}
		})
	}
}

func TestSpecMainnet_Phase0(t *testing.T) {
	testFork(t, phase0Fork)
}

func TestSpecMainnet_Altair(t *testing.T) {
	testFork(t, altairFork)
}

func TestSpecMainnet_Bellatrix(t *testing.T) {
	testFork(t, bellatrixFork)
}

func TestSpecMainnet_Capella(t *testing.T) {
	testFork(t, capellaFork)
}

func formatSpecFailure(errHeader, specFile, structName string, err error) string {
	return fmt.Sprintf("%s spec file=%s, struct=%s, err=%v",
		errHeader, specFile, structName, err)
}

func checkSSZEncoding(t *testing.T, f fork, fileName, structName string, base testCallback) {
	obj := base(f)
	if obj == nil {
		// skip
		return
	}
	output := readValidGenericSSZ(t, fileName, &obj)

	// Marshal
	res, err := obj.MarshalSSZTo(nil)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(res, output.ssz) {
		t.Fatal("bad marshalling")
	}

	// Unmarshal
	obj2 := base(f)
	if err := obj2.UnmarshalSSZ(res); err != nil {
		t.Fatal(formatSpecFailure("UnmarshalSSZ error", fileName, structName, err))
	}
	if !deepEqual(obj, obj2) {
		t.Fatal("bad unmarshalling")
	}

	// Root
	root, err := obj.HashTreeRoot()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(root[:], output.root) {
		fmt.Printf("%s bad root\n", fileName)
	}

	if objt, ok := obj.(codecTree); ok {
		// node root
		node, err := objt.GetTree()
		if err != nil {
			t.Fatal(err)
		}

		xx := node.Hash()
		if !bytes.Equal(xx, root[:]) {
			t.Fatal("bad node")
		}
	}
}

const (
	testsPath      = "./eth2.0-spec-tests/tests"
	serializedFile = "serialized.ssz_snappy"
	valueFile      = "value.yaml"
	rootsFile      = "roots.yaml"
)

func readDir(t *testing.T, path string) []string {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		t.Fatal(err)
	}
	res := []string{}
	for _, f := range files {
		res = append(res, filepath.Join(path, f.Name()))
	}
	return res
}

type output struct {
	root []byte
	ssz  []byte
}

func readValidGenericSSZ(t *testing.T, path string, obj interface{}) *output {
	serializedSnappy, err := ioutil.ReadFile(filepath.Join(path, serializedFile))
	if err != nil {
		t.Fatal(err)
	}
	serialized, err := snappy.Decode(nil, serializedSnappy)
	if err != nil {
		t.Fatal(err)
	}

	raw, err := ioutil.ReadFile(filepath.Join(path, valueFile))
	if err != nil {
		t.Fatal(err)
	}
	raw2, err := ioutil.ReadFile(filepath.Join(path, rootsFile))
	if err != nil {
		t.Fatal(err)
	}

	// Decode ssz root
	var out map[string]string
	if err := yaml.Unmarshal(raw2, &out); err != nil {
		t.Fatal(err)
	}
	root, err := hex.DecodeString(out["root"][2:])
	if err != nil {
		t.Fatal(err)
	}

	if err := ssz.UnmarshalSSZTest(raw, obj); err != nil {
		t.Fatal(err)
	}
	return &output{root: root, ssz: serialized}
}
