package main

import (
	"context"
	"github.com/Kotodian/registry"
	"github.com/Kotodian/registry/example/common"
	v1 "github.com/Kotodian/registry/pb/v1"
	ac_ocpp "github.com/Kotodian/registry/service/ac-ocpp"
	"google.golang.org/grpc"
	"log"
	"net"
)

type MasterServer struct {
	master *registry.Master
}

func (m *MasterServer) AddMember(ctx context.Context, req *v1.AddMemberReq) (*v1.AddMemberResp, error) {
	err := m.master.AddMember(ac_ocpp.NewSimpleService(req.GetHostname()))
	if err != nil {
		return nil, err
	}
	return &v1.AddMemberResp{}, nil
}

func main() {
	master, err := registry.NewMaster(common.RedisClient, nil, &ac_ocpp.AcOCPP{})
	server := grpc.NewServer()
	v1.RegisterMasterServer(server, &MasterServer{master: master})
	listener, err := net.Listen("tcp", ":8090")
	if err != nil {
		panic(err)
	}

	if err = server.Serve(listener); err != nil {
		panic(err)
	} else {
		log.Println("server listen: 8090")
	}
}
