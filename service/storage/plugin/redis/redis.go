package redis

import (
	"context"
	"github.com/Kotodian/registry/service"
	"github.com/Kotodian/registry/service/storage"
	"github.com/go-redis/redis/v8"
	"strings"
	"time"
)

type Storage struct {
	client redis.UniversalClient
}

func New(client redis.UniversalClient) storage.Storage {
	return &Storage{client: client}
}

func (s *Storage) Add(ctx context.Context, service service.SimpleService) error {
	if err := s.client.HSet(ctx, service.Key(), "hostname", strings.TrimPrefix(service.Key(), service.Prefix())).Err(); err != nil {
		return err
	}
	err := s.client.Expire(ctx, service.Key(), 5*time.Second).Err()
	if err != nil {
		return nil
	}
	return nil
}

func (s *Storage) Get(ctx context.Context, key string) (service.SimpleService, error) {
	panic("implement me")
}

func (s *Storage) Delete(ctx context.Context, key string) error {
	panic("implement me")
}

func (s *Storage) Exists(ctx context.Context, key string) (bool, error) {
	result, err := s.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	if result == 0 {
		return false, nil
	}
	return true, nil
}

func (s *Storage) KeepAlive(ctx context.Context, key string, duration time.Duration, stop chan struct{}) {
	ticker := time.NewTicker(duration)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-stop:
			return
		case <-ticker.C:
			err := s.client.Expire(ctx, key, 5*time.Second).Err()
			if err != nil {
				return
			}
		}
	}
}

func (s *Storage) Keys(ctx context.Context, prefix string) ([]string, error) {
	return s.client.Keys(ctx, prefix+"*").Result()

}
