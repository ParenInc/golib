package nats

import "github.com/pareninc/golib/logger"

// Existing callers pass *logger.Logger to NewClient; keep it satisfying Logger.
var _ Logger = (*logger.Logger)(nil)
