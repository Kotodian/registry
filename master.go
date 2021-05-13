package registry

import (
	"context"
	"github.com/Kotodian/registry/service"
	"github.com/go-redis/redis/v8"
	"log"
	"reflect"
	"strings"
	"time"
)

type Master struct {
	prefix      string
	members     *service.ServiceMap
	redisClient redis.UniversalClient
	kind        reflect.Type
	debug       bool
}

func NewMaster(redisClient redis.UniversalClient, svc service.SimpleService, debug bool) (*Master, error) {
	master := &Master{
		members:     service.NewServiceMap(),
		redisClient: redisClient,
		debug:       debug,
	}

	master.prefix = svc.Prefix()
	master.kind = reflect.TypeOf(svc)
	// 断掉重启时将redis信息重新同步到内存中
	err := master.start()
	if err != nil {
		return nil, err
	}
	go master.sync()
	return master, nil
}
func (m *Master) AddMember(worker service.SimpleService) error {
	result, err := m.redisClient.HGetAll(context.Background(),
		worker.Key()).
		Result()
	if err != nil {
		return err
	}
	if result == nil {
		m.members.Set(worker.Key(), worker)
	}

	log.Printf("service: %s register\n", worker.Key())
	return nil
}

func (m *Master) start() error {
	var results []string
	var err error
	results, err = m.redisClient.Keys(context.Background(), m.prefix+"*").Result()
	if err != nil {
		return err
	}
	if len(results) == 0 {
		if m.debug {
			log.Println("no service need to be reRegistered")
		}
		return nil
	}
	for _, result := range results {
		simpleService := service.Kind[m.kind](result)
		// todo: maps := m.redisClient.HValues(context.Background, result)
		// 		simpleService.Set(maps)
		m.members.Set(strings.TrimPrefix(result, m.prefix), simpleService)
		if m.debug {
			log.Printf("service %s registered.\n", result)
		}
	}
	return nil
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
		err := m.redisClient.HGetAll(context.Background(), key).Err()
		if err != nil {
			if err == redis.Nil {
				m.members.Delete(key)
				if m.debug {
					log.Printf("service %s unregister\n", key)
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
