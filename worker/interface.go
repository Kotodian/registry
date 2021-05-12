package worker

import (
	"context"
	"time"
)

type Worker interface {
	Prefix() string
	Key() string
	Register(ctx context.Context) error
	Heartbeat(ctx context.Context, duration time.Duration) error
	NotifyMaster(ctx context.Context) error
}
