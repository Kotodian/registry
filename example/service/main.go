package main

import (
	"context"
	"github.com/Kotodian/registry/example/common"
	v1 "github.com/Kotodian/registry/pb/v1"
	ac_ocpp2 "github.com/Kotodian/registry/service/plugin/ac-ocpp"
	"github.com/Kotodian/registry/service/storage/plugin/redis"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"log"
	"time"
)

func failOnErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
func main() {
	for i := 0; i < 10; i++ {
		go newWorker()
		time.Sleep(100 * time.Millisecond)
	}
	newWorker()
	select {}
}

func newWorker() {
	conn, err := grpc.Dial("127.0.0.1:8090", grpc.WithInsecure())
	failOnErr(err)

	ctx := context.Background()
	worker := ac_ocpp2.NewService(
		uuid.New().String(),
		redis.New(common.RedisClient),
		v1.NewMasterClient(conn))

	err = worker.Register(ctx)
	failOnErr(err)

	err = worker.NotifyMaster(ctx)
	failOnErr(err)

	err = worker.Heartbeat(ctx, 3*time.Second)
	failOnErr(err)
}
