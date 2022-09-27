package http

import (
	"encoding/hex"
	"fmt"

	consensus "github.com/umbracle/go-eth-consensus"
)

type ValidatorEndpoint struct {
	c *Client
}

func (c *Client) Validator() *ValidatorEndpoint {
	return &ValidatorEndpoint{c: c}
}

type AttesterDuty struct {
	PubKey                  string `json:"pubkey"`
	ValidatorIndex          uint   `json:"validator_index"`
	Slot                    uint64 `json:"slot"`
	CommitteeIndex          uint64 `json:"committee_index"`
	CommitteeLength         uint64 `json:"committee_length"`
	CommitteeAtSlot         uint64 `json:"committees_at_slot"`
	ValidatorCommitteeIndex uint64 `json:"validator_committee_index"`
}

func (v *ValidatorEndpoint) GetAttesterDuties(epoch uint64, indexes []string) ([]*AttesterDuty, error) {
	var out []*AttesterDuty
	err := v.c.Post(fmt.Sprintf("/eth/v1/validator/duties/attester/%d", epoch), indexes, &out)
	return out, err
}

type ProposerDuty struct {
	PubKey         string `json:"pubkey"`
	ValidatorIndex uint   `json:"validator_index"`
	Slot           uint64 `json:"slot"`
}

func (v *ValidatorEndpoint) GetProposerDuties(epoch uint64) ([]*ProposerDuty, error) {
	var out []*ProposerDuty
	err := v.c.Get(fmt.Sprintf("/eth/v1/validator/duties/proposer/%d", epoch), &out)
	return out, err
}

type CommitteeSyncDuty struct {
	PubKey                        string   `json:"pubkey"`
	ValidatorIndex                uint     `json:"validator_index"`
	ValidatorSyncCommitteeIndices []string `json:"validator_sync_committee_indices"`
}

func (v *ValidatorEndpoint) GetCommitteeSyncDuties(epoch uint64, indexes []string) ([]*CommitteeSyncDuty, error) {
	var out []*CommitteeSyncDuty
	err := v.c.Post(fmt.Sprintf("/eth/v1/validator/duties/sync/%d", epoch), indexes, &out)
	return out, err
}

func (v *ValidatorEndpoint) GetBlock(out consensus.BeaconBlock, slot uint64, randao [96]byte) error {
	buf := "0x" + hex.EncodeToString(randao[:])
	err := v.c.Get(fmt.Sprintf("/eth/v2/validator/blocks/%d?randao_reveal=%s", slot, buf), &out)
	return err
}

func (v *ValidatorEndpoint) RequestAttestationData(slot uint64, committeeIndex uint64) (*consensus.AttestationData, error) {
	var out *consensus.AttestationData
	err := v.c.Get(fmt.Sprintf("/eth/v1/validator/attestation_data?slot=%d&committee_index=%d", slot, committeeIndex), &out)
	return out, err
}

func (v *ValidatorEndpoint) AggregateAttestation(slot uint64, root [32]byte) (*consensus.Attestation, error) {
	var out *consensus.Attestation
	err := v.c.Get(fmt.Sprintf("/eth/v1/validator/aggregate_attestation?slot=%d&attestation_data_root=0x%s", slot, hex.EncodeToString(root[:])), &out)
	return out, err
}

func (v *ValidatorEndpoint) PublishAggregateAndProof(data []*consensus.SignedAggregateAndProof) error {
	err := v.c.Post("/eth/v1/validator/aggregate_and_proofs", data, nil)
	return err
}

type BeaconCommitteeSubscription struct {
	ValidatorIndex   uint64 `json:"validator_index"`
	Slot             uint64 `json:"slot"`
	CommitteeIndex   uint64 `json:"committee_index"`
	CommitteesAtSlot uint64 `json:"committee_at_slot"`
	IsAggregator     bool   `json:"is_aggregator"`
}

func (v *ValidatorEndpoint) BeaconCommitteeSubscriptions(subs []*BeaconCommitteeSubscription) error {
	err := v.c.Post("/eth/v1/validator/beacon_committee_subscriptions", subs, nil)
	return err
}

type SyncCommitteeSubscription struct {
	ValidatorIndex       uint64   `json:"validator_index"`
	SyncCommitteeIndices []uint64 `json:"sync_committee_indices"`
	UntilEpoch           uint64   `json:"until_epoch"`
}

func (v *ValidatorEndpoint) SyncCommitteeSubscriptions(subs []*SyncCommitteeSubscription) error {
	err := v.c.Post("/eth/v1/validator/sync_committee_subscriptions", subs, nil)
	return err
}

// produces a sync committee contribution
func (v *ValidatorEndpoint) SyncCommitteeContribution(slot uint64, subCommitteeIndex uint64, root [32]byte) (*consensus.SyncCommitteeContribution, error) {
	var out *consensus.SyncCommitteeContribution
	err := v.c.Get(fmt.Sprintf("/eth/v1/validator/sync_committee_contribution?slot=%d&subcommittee_index=%d&beacon_block_root=0x%s", slot, subCommitteeIndex, hex.EncodeToString(root[:])), &out)
	return out, err
}

func (v *ValidatorEndpoint) SubmitSignedContributionAndProof(signedContribution []*consensus.SignedContributionAndProof) error {
	err := v.c.Post("/eth/v1/validator/contribution_and_proofs", signedContribution, nil)
	return err
}

type ProposalPreparation struct {
	ValidatorIndex uint64
	FeeRecipient   [20]byte
}

func (v *ValidatorEndpoint) PrepareBeaconProposer(input []*ProposalPreparation) error {
	err := v.c.Post("/eth/v1/validator/prepare_beacon_proposer", input, nil)
	return err
}

type RegisterValidatorRequest struct {
	FeeRecipient [20]byte `json:"fee_recipient" ssz-size:"20"`
	GasLimit     uint64   `json:"gas_limit,string"`
	Timestamp    uint64   `json:"timestamp,string"`
	Pubkey       [48]byte `json:"pubkey" ssz-size:"48"`
}

type SignedValidatorRegistration struct {
	Message   *RegisterValidatorRequest `json:"message"`
	Signature [96]byte                  `json:"signature" ssz-size:"96"`
}

func (v *ValidatorEndpoint) RegisterValidator(msg []*SignedValidatorRegistration) error {
	err := v.c.Post("/eth/v1/validator/register_validator", msg, nil)
	return err
}
