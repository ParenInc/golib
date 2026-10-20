package logger

import (
	"bytes"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestFromLogrus_KeepsTheEntryOutputAndFields(t *testing.T) {
	var out bytes.Buffer
	base := logrus.New()
	base.Out = &out
	base.Formatter = &logrus.TextFormatter{DisableTimestamp: true}

	FromLogrus(base.WithField("component", "nats")).Infof("connected to %s", "nats://localhost")

	assert.Equal(t, "level=info msg=\"connected to nats://localhost\" component=nats\n", out.String())
}
