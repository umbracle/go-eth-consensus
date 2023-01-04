package consensus

import (
	"encoding/hex"
	"fmt"
	"strings"
)

type Domain [4]byte

func (d *Domain) UnmarshalText(data []byte) error {
	domainStr := string(data)
	if !strings.HasPrefix(domainStr, "0x") {
		return fmt.Errorf("not prefixed")
	}
	buf, err := hex.DecodeString(domainStr[2:])
	if err != nil {
		return err
	}
	if len(buf) != 4 {
		return fmt.Errorf("bad size")
	}
	copy(d[:], buf)
	return nil
}

func ToBytes96(b []byte) (res [96]byte) {
	copy(res[:], b)
	return
}

func ToBytes32(b []byte) (res [32]byte) {
	copy(res[:], b)
	return
}

func ComputeDomain(domain Domain, forkVersion [4]byte, genesisValidatorsRoot Root) ([32]byte, error) {
	// compute_fork_data_root
	// this returns the 32byte fork data root for the ``current_version`` and ``genesis_validators_root``.
	// This is used primarily in signature domains to avoid collisions across forks/chains.
	forkData := ForkData{
		CurrentVersion:        forkVersion,
		GenesisValidatorsRoot: genesisValidatorsRoot,
	}
	forkRoot, err := forkData.HashTreeRoot()
	if err != nil {
		return [32]byte{}, err
	}
	return ToBytes32(append(domain[:], forkRoot[:28]...)), nil
}

type DomainType Domain

var (
	DomainBeaconProposerType          = Domain{0, 0, 0, 0}
	DomainBeaconAttesterType          = Domain{1, 0, 0, 0}
	DomainRandaomType                 = Domain{2, 0, 0, 0}
	DomainDepositType                 = Domain{3, 0, 0, 0}
	DomainVoluntaryExitType           = Domain{4, 0, 0, 0}
	DomainSelectionProofType          = Domain{5, 0, 0, 0}
	DomainAggregateAndProofType       = Domain{6, 0, 0, 0}
	DomainSyncCommitteeType           = Domain{7, 0, 0, 0}
	DomainSyncCommitteeSelectionProof = Domain{8, 0, 0, 0}
	DomainContributionAndProof        = Domain{9, 0, 0, 0}
)
