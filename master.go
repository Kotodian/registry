package registry

import (
	"context"
	"github.com/Kotodian/registry/service"
	"github.com/Kotodian/registry/service/storage"
	"log"
	"reflect"
	"strings"
	"time"
)

type Master struct {
	prefix  string
	members *service.ServiceMap
	store   storage.Storage
	kind    reflect.Type
	debug   bool
}

func NewMaster(store storage.Storage, svc service.SimpleService, debug bool) (*Master, error) {
	master := &Master{
		members: service.NewServiceMap(),
		store:   store,
		debug:   debug,
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
	exists, err := m.store.Exists(context.Background(),
		worker.Key())
	if err != nil {
		return err
	}
	if exists {
		m.members.Set(strings.TrimPrefix(worker.Key(), m.prefix), worker)
	}

	log.Printf("service: %s register\n", worker.Key())
	return nil
}

func (m *Master) start() error {
	var results []string
	var err error
	results, err = m.store.Keys(context.Background(), m.prefix+"*")
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
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.workerSync()
		}
	}
}

func (m *Master) workerSync() {
	keys := m.members.Keys()
	if keys == nil {
		if m.debug {
			log.Printf("no need to sync\n")
		}
		return
	}
	for _, key := range keys {
		exists, err := m.store.Exists(context.Background(), m.prefix+key)
		if err != nil {
			m.members.Delete(key)
			if m.debug {
				log.Printf("service %s unregister\n", key)
			}
		}
		if !exists {
			m.members.Delete(key)
			if m.debug {
				log.Printf("service %s unregister\n", key)
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
