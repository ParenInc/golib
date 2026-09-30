package nats

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

// RetryPolicy controls how synchronous publishes carrying an idempotency key
// are retried.
// Retries reuse the same key, so the stream drops any copy that was stored
// despite the timeout; keep the total retry duration well under the stream's
// Duplicates window.
type RetryPolicy struct {
	// MaxAttempts is the total number of attempts, including the first.
	// Values <= 1 disable retries.
	MaxAttempts int
	// AttemptTimeout bounds how long a single attempt waits for the PubAck.
	AttemptTimeout time.Duration
	// InitialBackoff is the wait before the first retry. It doubles on each
	// subsequent retry, up to MaxBackoff, with jitter.
	InitialBackoff time.Duration
	MaxBackoff     time.Duration
}

// DefaultRetryPolicy retries up to 3 attempts, waiting at most ~16s in total.
var DefaultRetryPolicy = RetryPolicy{
	MaxAttempts:    3,
	AttemptTimeout: 5 * time.Second,
	InitialBackoff: 250 * time.Millisecond,
	MaxBackoff:     2 * time.Second,
}

// WithRetryPolicy overrides DefaultRetryPolicy for publishes carrying an
// idempotency key.
func WithRetryPolicy(policy RetryPolicy) Option {
	return func(o *clientOptions) {
		o.retryPolicy = policy
	}
}

func (p RetryPolicy) attempts() int {
	return max(p.MaxAttempts, 1)
}

// backoff returns the wait after the given failed attempt (1-based): half of
// the exponential delay plus up to the same amount again as jitter.
func (p RetryPolicy) backoff(attempt int) time.Duration {
	d := p.InitialBackoff << (attempt - 1)
	if d <= 0 || (p.MaxBackoff > 0 && d > p.MaxBackoff) {
		d = p.MaxBackoff
	}
	if d <= 0 {
		return 0
	}
	half := d / 2
	return half + rand.N(half+1)
}

// isRetryable reports whether err means the publish may not have been
// acknowledged (as opposed to the server rejecting it).
func isRetryable(err error) bool {
	return errors.Is(err, context.DeadlineExceeded) ||
		errors.Is(err, nats.ErrTimeout) ||
		errors.Is(err, nats.ErrNoResponders) ||
		errors.Is(err, nats.ErrDisconnected) ||
		errors.Is(err, jetstream.ErrNoStreamResponse)
}

type retryFunc func(attempt int, err error)

// publishWithRetry calls publish until it succeeds, fails with a
// non-retryable error, runs out of attempts, or ctx is done.
func publishWithRetry(ctx context.Context, policy RetryPolicy, publish func(context.Context) (*jetstream.PubAck, error), onRetry retryFunc) (*jetstream.PubAck, error) {
	attempts := policy.attempts()
	for attempt := 1; ; attempt++ {
		attemptCtx, cancel := ctx, context.CancelFunc(func() {})
		if policy.AttemptTimeout > 0 {
			attemptCtx, cancel = context.WithTimeout(ctx, policy.AttemptTimeout)
		}
		ack, err := publish(attemptCtx)
		cancel()
		if err == nil {
			return ack, nil
		}
		if ctx.Err() != nil || !isRetryable(err) {
			return nil, err
		}
		if attempt >= attempts {
			return nil, fmt.Errorf("publish failed after %d attempts: %w", attempts, err)
		}

		onRetry(attempt, err)
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(policy.backoff(attempt)):
		}
	}
}
