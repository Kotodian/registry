package ac_ocpp

import (
	"context"
	v1 "github.com/Kotodian/registry/pb/v1"
	"github.com/Kotodian/registry/worker"
	"github.com/go-redis/redis/v8"
	"time"
)

type AcOCPP struct {
	prefix       string
	hostname     string
	redisClient  *redis.Client
	masterClient v1.MasterClient
}

func (a *AcOCPP) Prefix() string {
	return a.prefix
}
func (a *AcOCPP) Key() string {
	return a.prefix + a.hostname
}

func NewWorker(prefix,
	hostname string,
	client *redis.Client,
	masterClient v1.MasterClient) worker.Worker {
	return &AcOCPP{prefix, hostname, client, masterClient}
}

func NewSimpleWorker(hostname string) worker.SimpleWorker {
	return &AcOCPP{hostname: hostname}
}

func (a *AcOCPP) Register(ctx context.Context) error {
	err := a.redisClient.HSet(ctx, a.Key(),
		"hostname", a.hostname).Err()
	if err != nil {
		return err
	}
	return nil
}

func (a *AcOCPP) Heartbeat(ctx context.Context, duration time.Duration) error {
	go func(ctx context.Context) {
		ticker := time.NewTicker(duration)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				err := a.redisClient.Expire(ctx, a.Key(), 20*time.Second).Err()
				if err != nil {
					return
				}
			}
		}
	}(ctx)
	return nil
}

func (a *AcOCPP) NotifyMaster(ctx context.Context) error {
	_, err := a.masterClient.AddMember(ctx, &v1.AddMemberReq{
		Hostname: a.hostname,
	})
	if err != nil {
		return err
	}
	return nil
}
