package service

import (
	"context"
	"reflect"
	"time"
)

type Service interface {
	SimpleService
	Register(ctx context.Context) error
	Heartbeat(ctx context.Context, duration time.Duration) error
	NotifyMaster(ctx context.Context) error
}

type SimpleService interface {
	Prefix() string
	Key() string
	Set(map[string]string)
}

var (
	Kind = make(map[reflect.Type]func(key string) SimpleService)
)
