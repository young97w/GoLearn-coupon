package biz

import (
	"account/internal"
	"account/log"
	"account/proto/pb"
	"context"
	"fmt"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"google.golang.org/grpc"
	"testing"
)

var client pb.AccountServiceClient

func init() {
	target := fmt.Sprintf("consul://%s:%d/%s?wait=14s", internal.AppConf.ConsulConfig.Host, internal.AppConf.ConsulConfig.Port, internal.AppConf.AccountSrvConfig.SrvName)
	conn, err := grpc.Dial(
		target,
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		log.Logger.Panic("grpcDial failed")
	}
	client = pb.NewAccountServiceClient(conn)
}

func TestAccountServer_AddAccount(t *testing.T) {
	res, err := client.AddAccount(context.Background(), &pb.AddAccountReq{
		Mobile:     "181111444661",
		Password:   "yynb666",
		Nickname:   "young",
		Gender:     "male",
		IsEmployee: false,
		Role:       0,
	})
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(res)
}
