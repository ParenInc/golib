package nats

import (
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

// Logger is the subset of logging methods the client uses to report connection
// events. *logger.Logger from this module satisfies it.
type Logger interface {
	Infof(format string, args ...any)
	Warnf(format string, args ...any)
	Errorf(format string, args ...any)
}

type Client struct {
	conn        *nats.Conn
	js          jetstream.JetStream
	logger      Logger
	retryPolicy RetryPolicy
}

type clientOptions struct {
	jsOpts      []jetstream.JetStreamOpt
	natsOpts    []nats.Option
	retryPolicy RetryPolicy
}

// Option customizes the client created by NewClient.
type Option func(*clientOptions)

// WithJetStreamOptions passes options through to the underlying JetStream
// context (e.g. jetstream.WithPublishAsyncMaxPending).
func WithJetStreamOptions(opts ...jetstream.JetStreamOpt) Option {
	return func(o *clientOptions) {
		o.jsOpts = append(o.jsOpts, opts...)
	}
}

// WithConnectionName names the connection, as shown by the server monitoring
// (e.g. nats server report connections).
func WithConnectionName(name string) Option {
	return func(o *clientOptions) {
		o.natsOpts = append(o.natsOpts, nats.Name(name))
	}
}

// WithInboxPrefix replaces the _INBOX prefix of the reply subjects (JetStream
// API, KV and publish acknowledgements), so the server permissions can keep
// the replies of the service away from the other users.
func WithInboxPrefix(prefix string) Option {
	return func(o *clientOptions) {
		o.natsOpts = append(o.natsOpts, nats.CustomInboxPrefix(prefix))
	}
}

func NewClient(config Configuration, logger Logger, opts ...Option) (*Client, error) {
	options := clientOptions{retryPolicy: DefaultRetryPolicy}
	for _, opt := range opts {
		opt(&options)
	}

	reconnectDelay := time.Second

	natsOpts := []nats.Option{}
	natsOpts = append(natsOpts, nats.ReconnectWait(reconnectDelay))
	natsOpts = append(natsOpts, nats.MaxReconnects(-1))
	natsOpts = append(natsOpts, nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
		logger.Warnf("Disconnected: will attempt reconnects: %v", err)
	}))
	natsOpts = append(natsOpts, nats.ReconnectHandler(func(nc *nats.Conn) {
		logger.Infof("Reconnected to %s", nc.ConnectedUrlRedacted())
	}))
	natsOpts = append(natsOpts, nats.ClosedHandler(func(nc *nats.Conn) {
		logger.Errorf("Exiting, no servers available")
	}))

	if config.User != "" && config.Password != "" {
		natsOpts = append(natsOpts, nats.UserInfo(config.User, config.Password))
	}
	natsOpts = append(natsOpts, options.natsOpts...)

	// connect to nats server
	nc, err := nats.Connect(config.URL, natsOpts...)
	if err != nil {
		return nil, err
	}

	// create jetstream context from nats connection
	js, err := jetstream.New(nc, options.jsOpts...)
	if err != nil {
		nc.Close()
		return nil, err
	}

	return &Client{
		conn:        nc,
		js:          js,
		logger:      logger,
		retryPolicy: options.retryPolicy,
	}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

// Drain stops the subscriptions, lets the in-flight messages and publishes
// complete, then closes the connection. It returns without waiting for the
// connection to be closed.
func (c *Client) Drain() error {
	return c.conn.Drain()
}

// JetStream returns the JetStream context of the client, for the features
// the client does not wrap (e.g. key-value buckets).
func (c *Client) JetStream() jetstream.JetStream {
	return c.js
}
