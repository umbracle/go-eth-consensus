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

func ComputeDomain(domain Domain, forkVersion [4]byte, genesisValidatorsRoot Root) ([]byte, error) {
	// compute_fork_data_root
	// this returns the 32byte fork data root for the ``current_version`` and ``genesis_validators_root``.
	// This is used primarily in signature domains to avoid collisions across forks/chains.
	forkData := ForkData{
		CurrentVersion:        forkVersion,
		GenesisValidatorsRoot: genesisValidatorsRoot,
	}
	forkRoot, err := forkData.HashTreeRoot()
	if err != nil {
		return nil, err
	}

	return append(domain[:], forkRoot[:28]...), nil
}
