package nats

import (
	"context"
	"time"

	"github.com/nats-io/nats.go/jetstream"
)

type Consumer interface {
	CreateConsumer(stream string, durable string, subjects []string) (jetstream.Consumer, error)
}

func (c *Client) CreateConsumer(stream string, durable string, subjects []string) (jetstream.Consumer, error) {
	conf := jetstream.ConsumerConfig{
		AckPolicy: jetstream.AckExplicitPolicy,
		AckWait:   time.Minute,
	}

	if durable != "" {
		conf.Durable = durable
	}

	if len(subjects) > 0 {
		conf.FilterSubjects = subjects
	}

	return c.js.CreateOrUpdateConsumer(context.Background(), stream, conf)
}
