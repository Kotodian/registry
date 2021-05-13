package ac_ocpp

import (
	"bytes"
	"context"
	"encoding/gob"
	v1 "github.com/Kotodian/registry/pb/v1"
	"github.com/Kotodian/registry/service"
	"github.com/Kotodian/registry/service/storage"
	"reflect"
	"time"
)

const (
	Prefix = "ac-ocpp:"
)

type AcOCPP struct {
	hostname       string
	store          storage.Storage
	masterClient   v1.MasterClient
	isRedisCluster bool
	stop           chan struct{}
}

func init() {
	service.Kind[reflect.TypeOf(&AcOCPP{})] = NewSimpleService
}

func (a *AcOCPP) Prefix() string {
	return Prefix
}
func (a *AcOCPP) Key() string {
	return Prefix + a.hostname
}

func NewService(
	hostname string,
	storage storage.Storage,
	masterClient v1.MasterClient) service.Service {
	svc := &AcOCPP{
		hostname:     hostname,
		store:        storage,
		masterClient: masterClient,
		stop:         make(chan struct{}),
	}

	return svc
}

func NewSimpleService(hostname string) service.SimpleService {
	return &AcOCPP{hostname: hostname}
}

func (a *AcOCPP) Register(ctx context.Context) error {
	if exists, err := a.store.Exists(ctx, a.Key()); err != nil {
		return err
	} else {
		if !exists {
			err = a.store.Add(ctx, a)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (a *AcOCPP) Heartbeat(ctx context.Context, duration time.Duration) error {
	go a.store.KeepAlive(ctx, a.Key(), duration, a.stop)
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

func (a *AcOCPP) UnRegister(ctx context.Context) error {
	err := a.store.Delete(ctx, a.Key())
	if err != nil {
		return err
	}
	a.stop <- struct{}{}
	return nil
}

func (a *AcOCPP) SimpleService() service.SimpleService {
	simpleService := &AcOCPP{}
	err := deepcopy(a, simpleService)
	if err != nil {
		return nil
	}
	return simpleService
}

func deepcopy(src interface{}, dst interface{}) error {
	var buffer bytes.Buffer
	if err := gob.NewEncoder(&buffer).Encode(src); err != nil {
		return err
	}

	return gob.NewDecoder(&buffer).Decode(dst)
}
