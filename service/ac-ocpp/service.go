package ac_ocpp

import (
	"context"
	v1 "github.com/Kotodian/registry/pb/v1"
	"github.com/Kotodian/registry/service"
	"github.com/go-redis/redis/v8"
	"log"
	"reflect"
	"time"
)

type AcOCPP struct {
	prefix         string
	hostname       string
	redisClient    redis.UniversalClient
	masterClient   v1.MasterClient
	isRedisCluster bool
}

func init() {
	service.Kind[reflect.TypeOf(&AcOCPP{})] = NewSimpleService
}

func (a *AcOCPP) Prefix() string {
	return a.prefix
}
func (a *AcOCPP) Key() string {
	return a.prefix + a.hostname
}

func NewService(
	hostname string,
	redisClient redis.UniversalClient,
	masterClient v1.MasterClient) service.Service {
	svc := &AcOCPP{
		prefix:       "ac-ocpp:",
		hostname:     hostname,
		redisClient:  redisClient,
		masterClient: masterClient,
	}

	return svc
}

func NewSimpleService(hostname string) service.SimpleService {
	return &AcOCPP{hostname: hostname, prefix: "ac-ocpp:"}
}

func (a *AcOCPP) Register(ctx context.Context) error {
	if err := a.redisClient.HSet(ctx, a.Key(), "hostname", a.hostname).Err(); err != nil {
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
				log.Printf("%s heartbeat\n", a.Key())
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

func (a *AcOCPP) Set(map[string]string) {

}
