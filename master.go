package registry

import (
	"context"
	"github.com/Kotodian/registry/worker"
	"github.com/go-redis/redis/v8"
	"time"
)

type Master struct {
	members worker.WorkerMap
	client  *redis.Client
}

func NewMaster(client *redis.Client) *Master {
	return &Master{
		members: worker.WorkerMap{},
		client:  client,
	}
}
func (m *Master) AddMember(worker worker.Worker) error {
	result, err := m.client.HGetAll(context.Background(),
		worker.Key()).
		Result()
	if err != nil {
		return err
	}
	if result == nil {
		m.members.Set(worker.Key(), worker)
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
	for _, key := range keys {
		err := m.client.HGetAll(context.Background(), key).Err()
		if err != nil {
			if err == redis.Nil {
				m.members.Delete(key)
			}
		}
	}
}

func (m *Master) Members() []worker.Worker {
	return m.members.Workers()
}

func (m *Master) Member(key string) worker.Worker {
	return m.members.Get(key)
}
