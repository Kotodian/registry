package registry

import (
	"context"
	"github.com/Kotodian/registry/service"
	"github.com/go-redis/redis/v8"
	"time"
)

type Master struct {
	members            *service.ServiceMap
	redisClient        *redis.Client
	redisClusterClient *redis.ClusterClient
	isRedisCluster     bool
}

func NewMaster(redisClient *redis.Client, redisClusterClient *redis.ClusterClient) *Master {
	master := &Master{
		members: service.NewServiceMap(),
	}
	if redisClusterClient != nil {
		master.isRedisCluster = true
		master.redisClusterClient = redisClusterClient
	} else {
		master.redisClient = redisClient
	}

	go master.sync()
	return master
}
func (m *Master) AddMember(worker service.SimpleService) error {
	if m.isRedisCluster {
		result, err := m.redisClusterClient.HGetAll(context.Background(),
			worker.Key()).
			Result()
		if err != nil {
			return err
		}
		if result == nil {
			m.members.Set(worker.Key(), worker)
		}
	} else {
		result, err := m.redisClient.HGetAll(context.Background(),
			worker.Key()).
			Result()
		if err != nil {
			return err
		}
		if result == nil {
			m.members.Set(worker.Key(), worker)
		}
	}
	return nil
}

func (m *Master) RmMember(key string) {
	m.members.Delete(key)
}

func (m *Master) sync() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	go func() {
		for {
			select {
			case <-ticker.C:
				m.workerSync()
			}
		}
	}()
}

func (m *Master) workerSync() {
	keys := m.members.Keys()
	if keys == nil {
		return
	}
	if m.isRedisCluster {
		for _, key := range keys {
			err := m.redisClusterClient.HGetAll(context.Background(), key).Err()
			if err != nil {
				if err == redis.Nil {
					m.members.Delete(key)
				}
			}
		}
	} else {
		for _, key := range keys {
			err := m.redisClient.HGetAll(context.Background(), key).Err()
			if err != nil {
				if err == redis.Nil {
					m.members.Delete(key)
				}
			}
		}
	}
}

func (m *Master) Members() []service.SimpleService {
	return m.members.Workers()
}

func (m *Master) Member(key string) service.SimpleService {
	return m.members.Get(key)
}
