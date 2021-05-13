package main

import (
	"context"
	"github.com/Kotodian/registry/example/common"
	v1 "github.com/Kotodian/registry/pb/v1"
	ac_ocpp "github.com/Kotodian/registry/service/ac-ocpp"
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
	newWorker()
	select {}
}

func newWorker() {
	conn, err := grpc.Dial("127.0.0.1:8090", grpc.WithInsecure())
	failOnErr(err)

	ctx := context.Background()
	worker := ac_ocpp.NewService(
		uuid.New().String(),
		common.RedisClient,
		v1.NewMasterClient(conn))

	err = worker.Register(ctx)
	failOnErr(err)

	err = worker.NotifyMaster(ctx)
	failOnErr(err)

	err = worker.Heartbeat(ctx, 3*time.Second)
	failOnErr(err)
}
