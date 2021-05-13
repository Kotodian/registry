package service

import (
	"context"
	"time"
)

type Service interface {
	Prefix() string
	Key() string
	Register(ctx context.Context) error
	Heartbeat(ctx context.Context, duration time.Duration) error
	NotifyMaster(ctx context.Context) error
}

type SimpleService interface {
	Key() string
}
