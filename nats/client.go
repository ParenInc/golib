package nats

import (
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/pareninc/golib/logger"
)

type Client struct {
	conn *nats.Conn
	js   jetstream.JetStream
}

func NewClient(config Configuration, logger *logger.Logger) (*Client, error) {
	reconnectDelay := time.Second

	opts := []nats.Option{}
	opts = append(opts, nats.ReconnectWait(reconnectDelay))
	opts = append(opts, nats.MaxReconnects(-1))
	opts = append(opts, nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
		logger.WithError(err).Warn("Disconnected: will attempt reconnects")
	}))
	opts = append(opts, nats.ReconnectHandler(func(nc *nats.Conn) {
		logger.WithField("url", nc.ConnectedUrl()).Info("Reconnected")
	}))
	opts = append(opts, nats.ClosedHandler(func(nc *nats.Conn) {
		logger.Error("Exiting, no servers available")
	}))

	// connect to nats server
	nc, err := nats.Connect(config.URL, opts...)
	if err != nil {
		return nil, err
	}

	// create jetstream context from nats connection
	js, err := jetstream.New(nc)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn: nc,
		js:   js,
	}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}
