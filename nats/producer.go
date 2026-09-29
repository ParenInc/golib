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
	PublishWithContext(ctx context.Context, subject string, data []byte, opts ...PublishOption) (*jetstream.PubAck, error)
	PublishAsync(subject string, data []byte, opts ...PublishOption) (jetstream.PubAckFuture, error)
}

type publishOptions struct {
	idempotencyKey string
}

// PublishOption customizes a single publish.
type PublishOption func(*publishOptions)

// WithIdempotencyKey sends key as the Nats-Msg-Id header, so the stream drops
// re-publishes of the same message within its Duplicates window. It also
// enables retries on timeouts according to the client's RetryPolicy; publishes
// without a key are never retried, since a retry could store a duplicate.
func WithIdempotencyKey(key string) PublishOption {
	return func(o *publishOptions) {
		o.idempotencyKey = key
	}
}

func newPublishOptions(opts []PublishOption) publishOptions {
	o := publishOptions{}
	for _, opt := range opts {
		opt(&o)
	}
	return o
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
// bounded by ctx. With WithIdempotencyKey, timed-out attempts are retried.
func (c *Client) PublishWithContext(ctx context.Context, subject string, data []byte, opts ...PublishOption) (*jetstream.PubAck, error) {
	o := newPublishOptions(opts)
	if o.idempotencyKey == "" {
		return c.js.Publish(ctx, subject, data)
	}

	return publishWithRetry(ctx, c.retryPolicy, func(attemptCtx context.Context) (*jetstream.PubAck, error) {
		return c.js.Publish(attemptCtx, subject, data, jetstream.WithMsgID(o.idempotencyKey))
	}, c.logRetry(subject, o.idempotencyKey))
}

// PublishAsync enqueues a publish without waiting for the server's PubAck. The
// returned future resolves once the server acknowledges or rejects the message.
// With WithIdempotencyKey, timed-out attempts are republished before the future
// resolves. The number of in-flight publishes can be capped with
// WithJetStreamOptions(jetstream.WithPublishAsyncMaxPending(n)); enqueue errors
// such as a full queue are returned immediately and not retried.
func (c *Client) PublishAsync(subject string, data []byte, opts ...PublishOption) (jetstream.PubAckFuture, error) {
	o := newPublishOptions(opts)
	if o.idempotencyKey == "" {
		return c.js.PublishAsync(subject, data)
	}

	publish := func() (jetstream.PubAckFuture, error) {
		return c.js.PublishAsync(subject, data, jetstream.WithMsgID(o.idempotencyKey))
	}
	first, err := publish()
	if err != nil {
		return nil, err
	}
	return newRetryingFuture(first, c.retryPolicy, publish, c.logRetry(subject, o.idempotencyKey)), nil
}

func (c *Client) logRetry(subject, key string) retryFunc {
	return func(attempt int, err error) {
		c.logger.Warnf("Retrying NATS publish to %s (key %s, attempt %d/%d failed): %v",
			subject, key, attempt, c.retryPolicy.attempts(), err)
	}
}
