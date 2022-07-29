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

func (v *ValidatorEndpoint) GetBlock(slot uint64, randao []byte) (*consensus.BeaconBlock, error) {
	buf := "0x" + hex.EncodeToString(randao)

	var out *consensus.BeaconBlock
	err := v.c.Get(fmt.Sprintf("/eth/v1/validator/blocks/%d?randao_reveal=%s", slot, buf), &out)

	return out, err
}

func (v *ValidatorEndpoint) RequestAttestationData(slot uint64, committeeIndex uint64) (*consensus.AttestationData, error) {
	var out *consensus.AttestationData
	err := v.c.Get(fmt.Sprintf("/eth/v1/validator/attestation_data?slot=%d&committee_index=%d", slot, committeeIndex), &out)
	return out, err
}

func (v *ValidatorEndpoint) AggregateAttestation(slot uint64, root []byte) (*consensus.Attestation, error) {
	var out *consensus.Attestation
	err := v.c.Get(fmt.Sprintf("/eth/v1/validator/aggregate_attestation?slot=%d&attestation_data_root=0x%s", slot, hex.EncodeToString(root[:])), &out)
	return out, err
}

func (v *ValidatorEndpoint) PublishAggregateAndProof(data []*consensus.SignedAggregateAndProof) error {
	err := v.c.Post("/eth/v1/validator/aggregate_and_proofs", data, nil)
	return err
}

// produces a sync committee contribution
func (v *ValidatorEndpoint) SyncCommitteeContribution(slot uint64, subCommitteeIndex uint64, root []byte) (*consensus.SyncCommitteeContribution, error) {
	var out *consensus.SyncCommitteeContribution
	err := v.c.Get(fmt.Sprintf("/eth/v1/validator/sync_committee_contribution?slot=%d&subcommittee_index=%d&beacon_block_root=0x%s", slot, subCommitteeIndex, hex.EncodeToString(root[:])), &out)
	return out, err
}

func (v *ValidatorEndpoint) SubmitSignedContributionAndProof(signedContribution []*consensus.SignedContributionAndProof) error {
	err := v.c.Post("/eth/v1/validator/contribution_and_proofs", signedContribution, nil)
	return err
}
