package nats

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testPolicy = RetryPolicy{
	MaxAttempts:    3,
	AttemptTimeout: time.Second,
	InitialBackoff: time.Millisecond,
	MaxBackoff:     2 * time.Millisecond,
}

var errRejected = errors.New("nats: maximum payload exceeded")

func noopRetry(int, error) {}

// syncResults returns a publish func yielding the given errors in order, then
// success, and a pointer to the number of calls made.
func syncResults(errs ...error) (func(context.Context) (*jetstream.PubAck, error), *int) {
	calls := 0
	return func(context.Context) (*jetstream.PubAck, error) {
		calls++
		if calls <= len(errs) {
			return nil, errs[calls-1]
		}
		return &jetstream.PubAck{Stream: "test"}, nil
	}, &calls
}

func TestPublishWithRetry(t *testing.T) {
	t.Run("should retry retryable errors until success", func(t *testing.T) {
		publish, calls := syncResults(context.DeadlineExceeded, nats.ErrTimeout)
		var retried []int

		ack, err := publishWithRetry(context.Background(), testPolicy, publish, func(attempt int, _ error) {
			retried = append(retried, attempt)
		})

		require.NoError(t, err)
		assert.Equal(t, "test", ack.Stream)
		assert.Equal(t, 3, *calls)
		assert.Equal(t, []int{1, 2}, retried)
	})

	t.Run("should not retry non-retryable errors", func(t *testing.T) {
		publish, calls := syncResults(errRejected)

		_, err := publishWithRetry(context.Background(), testPolicy, publish, noopRetry)

		assert.ErrorIs(t, err, errRejected)
		assert.Equal(t, 1, *calls)
	})

	t.Run("should give up after MaxAttempts", func(t *testing.T) {
		publish, calls := syncResults(nats.ErrTimeout, nats.ErrTimeout, nats.ErrTimeout, nats.ErrTimeout)

		_, err := publishWithRetry(context.Background(), testPolicy, publish, noopRetry)

		assert.ErrorIs(t, err, nats.ErrTimeout)
		assert.Equal(t, 3, *calls)
	})

	t.Run("should stop when the caller's context is done", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		calls := 0
		publish := func(context.Context) (*jetstream.PubAck, error) {
			calls++
			cancel()
			return nil, context.Canceled
		}

		_, err := publishWithRetry(ctx, testPolicy, publish, noopRetry)

		assert.ErrorIs(t, err, context.Canceled)
		assert.Equal(t, 1, calls)
	})

	t.Run("should bound each attempt by AttemptTimeout", func(t *testing.T) {
		policy := testPolicy
		policy.MaxAttempts = 2
		policy.AttemptTimeout = 10 * time.Millisecond
		calls := 0
		publish := func(ctx context.Context) (*jetstream.PubAck, error) {
			calls++
			<-ctx.Done()
			return nil, ctx.Err()
		}

		_, err := publishWithRetry(context.Background(), policy, publish, noopRetry)

		assert.ErrorIs(t, err, context.DeadlineExceeded)
		assert.Equal(t, 2, calls)
	})
}

func TestRetryPolicy_Backoff(t *testing.T) {
	policy := RetryPolicy{InitialBackoff: 100 * time.Millisecond, MaxBackoff: 300 * time.Millisecond}

	for range 100 {
		first := policy.backoff(1)
		assert.GreaterOrEqual(t, first, 50*time.Millisecond)
		assert.LessOrEqual(t, first, 100*time.Millisecond)

		capped := policy.backoff(10)
		assert.GreaterOrEqual(t, capped, 150*time.Millisecond)
		assert.LessOrEqual(t, capped, 300*time.Millisecond)
	}

	assert.Equal(t, time.Duration(0), RetryPolicy{}.backoff(1))
}
