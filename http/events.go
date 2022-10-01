package http

import (
	"context"
	"fmt"
	"strings"

	"github.com/r3labs/sse"
	consensus "github.com/umbracle/go-eth-consensus"
)

type HeadEvent struct {
	Slot                      uint64   `json:"slot"`
	Block                     [32]byte `json:"block"`
	State                     [32]byte `json:"state"`
	EpochTransition           bool     `json:"epoch_transition"`
	CurrentDutyDependentRoot  [32]byte `json:"current_duty_dependent_root"`
	PreviousDutyDependentRoot [32]byte `json:"previous_duty_dependent_root"`
	ExecutionOptimistic       bool     `json:"execution_optimistic"`
}

type BlockEvent struct {
	Slot                uint64   `json:"slot"`
	Block               [32]byte `json:"block"`
	ExecutionOptimistic bool     `json:"execution_optimistic"`
}

type FinalizedCheckpointEvent struct {
	Block               [32]byte `json:"block"`
	State               [32]byte `json:"state"`
	Epoch               uint64   `json:"epoch"`
	ExecutionOptimistic bool     `json:"execution_optimistic"`
}

type ChainReorgEvent struct {
	Slot                uint64   `json:"slot"`
	Depth               uint64   `json:"depth"`
	OldHeadBlock        [32]byte `json:"old_head_block"`
	NewHeadBlock        [32]byte `json:"new_head_block"`
	OldHeadState        [32]byte `json:"old_head_state"`
	NewHeadState        [32]byte `json:"new_head_state"`
	Epoch               uint64   `json:"epoch"`
	ExecutionOptimistic bool     `json:"execution_optimistic"`
}

var eventValidTopics = []string{
	"head", "block", "attestation", "finalized_checkpoint",
}

func isValidTopic(str string) bool {
	for _, topic := range eventValidTopics {
		if str == topic {
			return true
		}
	}
	return false
}

func (c *Client) Events(ctx context.Context, topics []string, handler func(obj interface{})) error {
	for _, topic := range topics {
		if !isValidTopic(topic) {
			return fmt.Errorf("topic '%s' is not valid", topic)
		}
	}

	client := sse.NewClient(c.url + "/eth/v1/events?topics=" + strings.Join(topics, ","))
	if err := client.SubscribeRawWithContext(ctx, func(msg *sse.Event) {
		var obj interface{}
		switch string(msg.Event) {
		case "head":
			obj = new(HeadEvent)
		case "block":
			obj = new(BlockEvent)
		case "attestation":
			obj = new(consensus.Attestation)
		case "voluntary_exit":
			obj = new(consensus.SignedVoluntaryExit)
		case "finalized_checkpoint":
			obj = new(FinalizedCheckpointEvent)
		case "chain_reorg":
			obj = new(ChainReorgEvent)
		case "contribution_and_proof":
			obj = new(consensus.SignedContributionAndProof)
		default:
			c.config.logger.Printf("[DEBUG]: event not tracked: %s", string(msg.Event))
			return
		}

		if err := Unmarshal(msg.Data, obj, c.config.untrackedKeys); err != nil {
			c.config.logger.Printf("[ERROR]: failed to decode %s event: %v", string(msg.Event), err)
			return
		}
		handler(obj)
	}); err != nil {
		return err
	}
	return nil
}
