package ac_ocpp

import (
	"context"
	v1 "github.com/Kotodian/registry/pb/v1"
	"github.com/Kotodian/registry/service"
	"github.com/go-redis/redis/v8"
	"time"
)

type AcOCPP struct {
	prefix             string
	hostname           string
	redisClient        *redis.Client
	redisClusterClient *redis.ClusterClient
	masterClient       v1.MasterClient
	isRedisCluster     bool
}

func (a *AcOCPP) Prefix() string {
	return a.prefix
}
func (a *AcOCPP) Key() string {
	return a.prefix + a.hostname
}

func NewService(prefix,
	hostname string,
	redisClient *redis.Client,
	redisClusterClient *redis.ClusterClient,
	masterClient v1.MasterClient) service.Service {
	svc := &AcOCPP{
		prefix:       prefix,
		hostname:     hostname,
		masterClient: masterClient,
	}
	if redisClusterClient != nil {
		svc.redisClusterClient = redisClusterClient
		svc.isRedisCluster = true
	} else {
		svc.redisClient = redisClient
	}
	return svc
}

func NewSimpleService(hostname string) service.Service {
	return &AcOCPP{hostname: hostname}
}

func (a *AcOCPP) Register(ctx context.Context) error {
	if a.isRedisCluster {
		if err := a.redisClusterClient.HSet(ctx, a.Key(), "hostname", a.hostname).Err(); err != nil {
			return err
		}
	} else {
		if err := a.redisClient.HSet(ctx, a.Key(), "hostname", a.hostname).Err(); err != nil {
			return err
		}
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
				if a.isRedisCluster {
					err := a.redisClusterClient.Expire(ctx, a.Key(), 20*time.Second).Err()
					if err != nil {
						return
					}
				} else {
					err := a.redisClient.Expire(ctx, a.Key(), 20*time.Second).Err()
					if err != nil {
						return
					}
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
