package nats

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go/jetstream"
)

type Producer interface {
	CreateOrUpdateStream(streamName string, subjects []string) (jetstream.Stream, error)
	EnsureStreamExists(streamName string) error
	Publish(subject string, data []byte) error
}

func (c *Client) CreateOrUpdateStream(streamName string, subjects []string) (jetstream.Stream, error) {
	return c.js.CreateOrUpdateStream(context.Background(), jetstream.StreamConfig{
		Name:     streamName,
		Subjects: subjects,
	})
}

func (c *Client) EnsureStreamExists(streamName string) error {
	_, err := c.js.Stream(context.Background(), streamName)
	if err != nil {
		return fmt.Errorf("stream %q not found: %w", streamName, err)
	}
	return nil
}

func (c *Client) Publish(subject string, data []byte) error {
	_, err := c.js.Publish(context.Background(), subject, data)
	return err
}
