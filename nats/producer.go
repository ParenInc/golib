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
	PublishWithContext(ctx context.Context, subject string, data []byte, opts ...jetstream.PublishOpt) (*jetstream.PubAck, error)
	PublishAsync(subject string, data []byte, opts ...jetstream.PublishOpt) (jetstream.PubAckFuture, error)
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
	_, err := c.PublishWithContext(context.Background(), subject, data)
	return err
}

// PublishWithContext publishes synchronously and waits for the server's PubAck,
// bounded by ctx.
func (c *Client) PublishWithContext(ctx context.Context, subject string, data []byte, opts ...jetstream.PublishOpt) (*jetstream.PubAck, error) {
	return c.js.Publish(ctx, subject, data, opts...)
}

// PublishAsync enqueues a publish without waiting for the server's PubAck. The
// returned future resolves once the server acknowledges or rejects the message.
// The number of in-flight publishes can be capped with
// WithJetStreamOptions(jetstream.WithPublishAsyncMaxPending(n)).
func (c *Client) PublishAsync(subject string, data []byte, opts ...jetstream.PublishOpt) (jetstream.PubAckFuture, error) {
	return c.js.PublishAsync(subject, data, opts...)
}
