package nats

import (
	"context"
	"testing"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	natsserver "github.com/nats-io/nats-server/v2/test"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type nopLogger struct{}

func (nopLogger) Infof(string, ...any)  {}
func (nopLogger) Warnf(string, ...any)  {}
func (nopLogger) Errorf(string, ...any) {}

func runServer(t *testing.T) *server.Server {
	t.Helper()
	opts := natsserver.DefaultTestOptions
	opts.Port = -1
	opts.JetStream = true
	opts.StoreDir = t.TempDir()
	s := natsserver.RunServer(&opts)
	t.Cleanup(s.Shutdown)
	return s
}

func newTestClient(t *testing.T, s *server.Server, opts ...Option) *Client {
	t.Helper()
	c, err := NewClient(Configuration{URL: s.ClientURL()}, nopLogger{}, opts...)
	require.NoError(t, err)
	t.Cleanup(c.Close)
	return c
}

func TestNewClient_ConnectionOptions(t *testing.T) {
	s := runServer(t)
	c := newTestClient(t, s, WithConnectionName("golib-test"), WithInboxPrefix("_INBOX_golib"))

	assert.Equal(t, "golib-test", c.conn.Opts.Name)
	assert.Contains(t, c.conn.NewRespInbox(), "_INBOX_golib.")

	// the JetStream API replies use the custom inbox prefix
	observer, err := nats.Connect(s.ClientURL())
	require.NoError(t, err)
	t.Cleanup(observer.Close)
	replies, err := observer.SubscribeSync("_INBOX_golib.>")
	require.NoError(t, err)
	require.NoError(t, observer.Flush())

	_, err = c.JetStream().AccountInfo(context.Background())
	require.NoError(t, err)
	_, err = replies.NextMsg(time.Second)
	assert.NoError(t, err, "the reply must be sent to the custom inbox")
}

func TestClient_JetStreamGivesAccessToKeyValue(t *testing.T) {
	s := runServer(t)
	c := newTestClient(t, s)
	ctx := context.Background()

	kv, err := c.JetStream().CreateKeyValue(ctx, jetstream.KeyValueConfig{Bucket: "status"})
	require.NoError(t, err)
	_, err = kv.Put(ctx, "k", []byte("v"))
	require.NoError(t, err)
	entry, err := kv.Get(ctx, "k")
	require.NoError(t, err)
	assert.Equal(t, []byte("v"), entry.Value())
}

func TestClient_PublishWithExpectStream(t *testing.T) {
	s := runServer(t)
	c := newTestClient(t, s)
	ctx := context.Background()
	_, err := c.CreateOrUpdateStream("REPORTS", []string{"report.>"})
	require.NoError(t, err)

	tests := map[string]struct {
		opts    []PublishOption
		wantErr bool
	}{
		"when the subject is stored in the expected stream": {
			opts: []PublishOption{WithExpectStream("REPORTS")},
		},
		"when the subject is stored in another stream": {
			opts:    []PublishOption{WithExpectStream("OTHER")},
			wantErr: true,
		},
		"with an idempotency key": {
			opts: []PublishOption{WithExpectStream("REPORTS"), WithIdempotencyKey("id-1")},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ack, err := c.PublishWithContext(ctx, "report.requested", []byte("{}"), tt.opts...)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, "REPORTS", ack.Stream)

			future, err := c.PublishAsync("report.requested", []byte("{}"), tt.opts...)
			require.NoError(t, err)
			select {
			case ack := <-future.Ok():
				assert.Equal(t, "REPORTS", ack.Stream)
			case err := <-future.Err():
				t.Fatalf("async publish failed: %v", err)
			case <-time.After(5 * time.Second):
				t.Fatal("no async publish acknowledgement")
			}
		})
	}
}

func TestClient_DrainClosesTheConnection(t *testing.T) {
	s := runServer(t)
	c := newTestClient(t, s)

	require.NoError(t, c.Drain())

	assert.Eventually(t, c.conn.IsClosed, 5*time.Second, 10*time.Millisecond)
}
