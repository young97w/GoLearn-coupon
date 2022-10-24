package biz

import (
	"context"
	"coupon/internal"
	"coupon/log"
	"coupon/proto/pb"
	"fmt"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"google.golang.org/grpc"
	"testing"
)

var client pb.CouponServiceClient

func init() {
	target := fmt.Sprintf("consul://%s:%d/%s?wait=14s", internal.AppConf.ConsulConfig.Host, internal.AppConf.ConsulConfig.Port, internal.AppConf.CouponSrvConfig.SrvName)
	conn, err := grpc.Dial(
		target,
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
	)
	if err != nil {
		log.Logger.Panic("grpcDial failed" + err.Error())
	}
	client = pb.NewCouponServiceClient(conn)
}

func TestCouponServer_AddCoffee(t *testing.T) {
	for i := 0; i < 2; i++ {
		_, err := client.AddCoffee(context.Background(), &pb.CoffeeItem{
			Name:        fmt.Sprintf("红茶%d", i),
			Price:       29.9,
			RealPrice:   19.9,
			SoldNum:     0,
			Sku:         "BT-01",
			Description: "black tea",
			Image:       "",
		})
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestCouponServer_DeleteCoffee(t *testing.T) {
	ids := []int32{41, 42}
	_, err := client.DeleteCoffee(context.Background(), &pb.DeleteCoffeeReq{Id: ids})
	if err != nil {
		t.Fatal(err)
	}
}

func TestCouponServer_UpdateCoffee(t *testing.T) {
	_, err := client.UpdateCoffee(context.Background(), &pb.CoffeeItem{
		Id:          31,
		Name:        "升级过的咖啡",
		Price:       35.9,
		RealPrice:   22.9,
		SoldNum:     0,
		Sku:         "CF-UP",
		Description: "666",
		Image:       "",
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestCouponServer_SellCoffee(t *testing.T) {
	m.Lock()
	defer m.Unlock()
	_, err := client.SellCoffee(context.Background(), &pb.SellCoffeeReq{
		SoldNum: 2,
		Id:      31,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestCouponServer_ListCoffee(t *testing.T) {
	listRes, err := client.ListCoffee(context.Background(), &pb.ListCoffeeReq{
		PageSize: 10,
		PageNo:   1,
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(listRes)
}

func TestCouponServer_GetCoffee(t *testing.T) {
	coffee, err := client.GetCoffee(context.Background(), &pb.CoffeeItem{Id: 31})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(coffee)
}
