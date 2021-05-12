package ac_ocpp

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type AcOCPP struct {
	prefix   string
	hostname string
	client   *redis.Client
}

func (a *AcOCPP) Prefix() string {
	return a.prefix
}
func (a *AcOCPP) Key() string {
	return a.prefix + a.hostname
}

func NewWorker(prefix, hostname string, client *redis.Client) *AcOCPP {
	return &AcOCPP{prefix, hostname, client}
}

func (a *AcOCPP) Heartbeat(ctx context.Context, duration time.Duration) error {
	key := a.prefix + a.hostname
	err := a.client.HSet(ctx, key,
		"hostname", a.hostname).Err()
	if err != nil {
		return err
	}
	// todo: notify master
	go func(ctx context.Context) {
		ticker := time.NewTicker(duration)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				err = a.client.Expire(ctx, key, 20*time.Second).Err()
				if err != nil {
					return
				}
			}
		}
	}(ctx)
	return nil
}
