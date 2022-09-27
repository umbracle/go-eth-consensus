package http

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/r3labs/sse"
)

type HeadEvent struct {
	Slot                      string
	Block                     string
	State                     string
	EpochTransition           bool
	CurrentDutyDependentRoot  string
	PreviousDutyDependentRoot string
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
		switch string(msg.Event) {
		case "head":
			var headEvent *HeadEvent
			if err := json.Unmarshal(msg.Data, &headEvent); err != nil {
				c.config.logger.Printf("[ERROR]: failed to decode head event: %v", err)
			} else {
				handler(err)
			}

		default:
			c.config.logger.Printf("[DEBUG]: event not tracked: %s", string(msg.Event))
		}
	}); err != nil {
		return err
	}
	return nil
}
