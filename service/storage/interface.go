package storage

import (
	"context"
	"github.com/Kotodian/registry/service"
	"time"
)

type Storage interface {
	Keys(ctx context.Context, prefix string) ([]string, error)
	Add(ctx context.Context, service service.SimpleService) error
	Get(ctx context.Context, key string) (service.SimpleService, error)
	KeepAlive(ctx context.Context, key string, duration time.Duration, stop chan struct{})
	Exists(ctx context.Context, key string) (bool, error)
	Delete(ctx context.Context, key string) error
}
