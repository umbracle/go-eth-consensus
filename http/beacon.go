package http

import (
	"fmt"

	consensus "github.com/umbracle/go-eth-consensus"
)

type BeaconEndpoint struct {
	c *Client
}

func (c *Client) Beacon() *BeaconEndpoint {
	return &BeaconEndpoint{c: c}
}

type GenesisInfo struct {
	Time uint64   `json:"genesis_time"`
	Root [32]byte `json:"genesis_validators_root"`
	Fork string   `json:"genesis_fork_version"`
}

func (b *BeaconEndpoint) Genesis() (*GenesisInfo, error) {
	var out GenesisInfo
	err := b.c.Get("/eth/v1/beacon/genesis", &out)
	return &out, err
}

func (b *BeaconEndpoint) SubmitCommitteeDuties(duties []*consensus.SyncCommitteeMessage) error {
	err := b.c.Post("/eth/v1/beacon/pool/sync_committees", duties, nil)
	return err
}

type Validator struct {
	Index     uint64             `json:"index"`
	Balance   uint64             `json:"balance"`
	Status    string             `json:"status"`
	Validator *ValidatorMetadata `json:"validator"`
}

type ValidatorMetadata struct {
	PubKey                     string `json:"pubkey"`
	WithdrawalCredentials      string `json:"withdrawal_credentials"`
	EffectiveBalance           uint64 `json:"effective_balance"`
	Slashed                    bool   `json:"slashed"`
	ActivationElegibilityEpoch uint64 `json:"activation_eligibility_epoch"`
	ActivationEpoch            uint64 `json:"activation_epoch"`
	ExitEpoch                  uint64 `json:"exit_epoch"`
	WithdrawableEpoch          uint64 `json:"withdrawable_epoch"`
}

type StateId interface {
	StateID() string
}

type BlockId interface {
	BlockID() string
}

type plainStateId string

func (s plainStateId) StateID() string {
	return string(s)
}

func (s plainStateId) BlockID() string {
	return s.StateID()
}

type Slot uint64

func (s Slot) StateID() string {
	return fmt.Sprintf("%d", s)
}

func (s Slot) BlockID() string {
	return s.StateID()
}

const (
	Head      plainStateId = "head"
	Genesis   plainStateId = "genesis"
	Finalized plainStateId = "finalized"
)

func (b *BeaconEndpoint) GetValidators(id StateId) ([]*Validator, error) {
	var out []*Validator
	err := b.c.Get("/eth/v1/beacon/states/"+id.StateID()+"/validators", &out)
	return out, err
}

func (b *BeaconEndpoint) GetValidatorByPubKey(pub string, id StateId) (*Validator, error) {
	var out *Validator
	err := b.c.Get("/eth/v1/beacon/states/"+id.StateID()+"/validators/"+pub, &out)
	return out, err
}

func (b *BeaconEndpoint) PublishSignedBlock(block consensus.SignedBeaconBlock) error {
	err := b.c.Post("/eth/v1/beacon/blocks", block, nil)
	return err
}

func (b *BeaconEndpoint) PublishAttestations(data []*consensus.Attestation) error {
	err := b.c.Post("/eth/v1/beacon/pool/attestations", data, nil)
	return err
}

type Block struct {
	Message   consensus.BeaconBlock
	Signature [96]byte
}

func (b *BeaconEndpoint) GetBlock(id BlockId, block consensus.BeaconBlock) (*Block, error) {
	out := &Block{
		Message: block,
	}
	err := b.c.Get("/eth/v2/beacon/blocks/"+id.BlockID(), out)
	return out, err
}

type BlockHeaderResponse struct {
	Root      [32]byte     `json:"root"`
	Canonical bool         `json:"canonical"`
	Header    *BlockHeader `json:"header"`
}

type BlockHeader struct {
	Message   *consensus.BeaconBlockHeader `json:"message"`
	Signature [96]byte                     `json:"signature"`
}

func (b *BeaconEndpoint) GetBlockHeader(id BlockId) (*BlockHeaderResponse, error) {
	var out *BlockHeaderResponse
	err := b.c.Get("/eth/v1/beacon/headers/"+id.BlockID(), &out)
	return out, err
}

func (b *BeaconEndpoint) GetBlockRoot(id BlockId) ([32]byte, error) {
	var data struct {
		Root [32]byte
	}
	err := b.c.Get("/eth/v1/beacon/blocks/"+id.BlockID()+"/root", &data)
	return data.Root, err
}

func (b *BeaconEndpoint) GetBlockAttestations(id BlockId) ([]*consensus.Attestation, error) {
	var out []*consensus.Attestation
	err := b.c.Get("/eth/v1/beacon/blocks/"+id.BlockID()+"/attestations", &out)
	return out, err
}
