package worker

import (
	"context"
	"time"
)

type Worker interface {
	Prefix() string
	Key() string
	Heartbeat(ctx context.Context, duration time.Duration) error
}
