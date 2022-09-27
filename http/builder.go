package http

import (
	"fmt"

	consensus "github.com/umbracle/go-eth-consensus"
)

type BuilderEndpoint struct {
	c *Client
}

func (c *Client) Builder() *BuilderEndpoint {
	return &BuilderEndpoint{c: c}
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

func (b *BuilderEndpoint) RegisterValidator(msg []*SignedValidatorRegistration) error {
	err := b.c.Post("/eth/v1/builder/validators", msg, nil)
	return err
}

type BuilderBid struct {
	Header *consensus.ExecutionPayloadHeader `json:"header"`
	Value  consensus.Uint256                 `json:"value" ssz-size:"32"`
	Pubkey [48]byte                          `json:"pubkey" ssz-size:"48"`
}

type SignedBuilderBid struct {
	Message   *BuilderBid `json:"message"`
	Signature [96]byte    `json:"signature" ssz-size:"96"`
}

func (b *BuilderEndpoint) GetExecutionPayload(slot uint64, parentHash [32]byte, pubKey [48]byte) (*SignedBuilderBid, error) {
	var out *SignedBuilderBid
	err := b.c.Get(fmt.Sprintf("/eth/v1/builder/header/%d/0x%x/0x%x", slot, parentHash[:], pubKey[:]), &out)
	return out, err
}

func (b *BuilderEndpoint) SubmitBlindedBlock(msg *consensus.SignedBlindedBeaconBlock) (*consensus.ExecutionPayload, error) {
	var out *consensus.ExecutionPayload
	err := b.c.Post("/eth/v1/builder/blinded_blocks", msg, &out)
	return out, err
}

func (b *BuilderEndpoint) Status() (bool, error) {
	return b.c.Status("/eth/v1/builder/status")
}
