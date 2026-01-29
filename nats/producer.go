package nats

import (
	"context"

	"github.com/nats-io/nats.go/jetstream"
)

type Producer interface {
	CreateOrUpdateStream(streamName string, subjects []string) (jetstream.Stream, error)
	Publish(subject string, data []byte) error
}

func (c *Client) CreateOrUpdateStream(streamName string, subjects []string) (jetstream.Stream, error) {
	return c.js.CreateOrUpdateStream(context.Background(), jetstream.StreamConfig{
		Name:     streamName,
		Subjects: subjects,
	})
}

func (c *Client) Publish(subject string, data []byte) error {
	_, err := c.js.Publish(context.Background(), subject, data)
	return err
}
