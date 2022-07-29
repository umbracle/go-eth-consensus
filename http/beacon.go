package http

import (
	consensus "github.com/umbracle/go-eth-consensus"
)

type BeaconEndpoint struct {
	c *Client
}

func (c *Client) Beacon() *BeaconEndpoint {
	return &BeaconEndpoint{c: c}
}

type Genesis struct {
	Time uint64 `json:"genesis_time"`
	Root []byte `json:"genesis_validators_root"`
	Fork string `json:"genesis_fork_version"`
}

func (b *BeaconEndpoint) Genesis() (*Genesis, error) {
	var out Genesis
	err := b.c.Get("/eth/v1/beacon/genesis", &out)
	return &out, err
}

type SyncCommitteeMessage struct {
	Slot           uint64 `json:"slot"`
	BlockRoot      []byte `json:"beacon_block_root"`
	ValidatorIndex uint64 `json:"validator_index"`
	Signature      []byte `json:"signature"`
}

func (b *BeaconEndpoint) SubmitCommitteeDuties(duties []*SyncCommitteeMessage) error {
	err := b.c.Post("/eth/v1/beacon/pool/sync_committees", duties, nil)
	return err
}

type Validator struct {
	Index     uint64             `json:"index"`
	Status    string             `json:"status"`
	Validator *ValidatorMetadata `json:"validator"`
}

type ValidatorMetadata struct {
	PubKey                     string `json:"pubkey"`
	Slashed                    bool   `json:"slashed"`
	ActivationElegibilityEpoch uint64 `json:"activation_eligibility_epoch"`
	ActivationEpoch            uint64 `json:"activation_epoch"`
	ExitEpoch                  uint64 `json:"exit_epoch"`
	WithdrawableEpoch          uint64 `json:"withdrawable_epoch"`
}

func (b *BeaconEndpoint) GetValidatorByPubKey(pub string) (*Validator, error) {
	var out *Validator
	err := b.c.Get("/eth/v1/beacon/states/head/validators/"+pub, &out)
	return out, err
}

func (b *BeaconEndpoint) PublishSignedBlock(block *consensus.SignedBeaconBlock) error {
	err := b.c.Post("/eth/v1/beacon/blocks", block, nil)
	return err
}

func (b *BeaconEndpoint) PublishAttestations(data []*consensus.Attestation) error {
	err := b.c.Post("/eth/v1/beacon/pool/attestations", data, nil)
	return err
}

func (b *BeaconEndpoint) GetHeadBlockRoot() ([]byte, error) {
	var data struct {
		Root []byte
	}
	err := b.c.Get("/eth/v1/beacon/blocks/head/root", &data)
	return data.Root, err
}
