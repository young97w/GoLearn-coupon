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
	//res, err := client.AddAccount(context.Background(), &pb.AddAccountReq{
	//	Mobile:     "181111444661",
	//	Password:   "yynb666",
	//	Nickname:   "young",
	//	Gender:     "male",
	//	IsEmployee: false,
	//	Role:       0,
	//})

	res, err := client.AddAccount(context.Background(), &pb.AddAccountReq{
		Mobile:     "181111444662",
		Password:   "yynb662",
		Nickname:   "young2",
		Gender:     "male",
		IsEmployee: true,
		Role:       0,
	})

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(res)
}

func TestAccountServer_UpdateAccount(t *testing.T) {
	res, err := client.UpdateAccount(context.Background(), &pb.UpdateAccountReq{
		Id:     4,
		Role:   2,
		Mobile: "18111444662",
	})
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(res)
}

func TestAccountServer_GetAccountList(t *testing.T) {
	res, err := client.GetAccountList(context.Background(), &pb.ListAccountReq{
		PageNo:   1,
		PageSize: 10,
	})
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(res)
}

func TestAccountServer_GetAccountByMobile(t *testing.T) {
	res, err := client.GetAccountByMobile(context.Background(), &pb.MobileAccountReq{Mobile: "18111444661"})
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(res)
}

func TestAccountServer_GetAccountById(t *testing.T) {
	res, err := client.GetAccountById(context.Background(), &pb.IdAccountReq{Id: 1})
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(res)
}
