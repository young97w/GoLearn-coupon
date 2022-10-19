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

func TestCouponServer_AddCoupon(t *testing.T) {
	res, err := client.AddCoupon(context.Background(), &pb.AddCouponReq{
		Amount:       5,
		Name:         "员工券-20起八折",
		CouponType:   2, //1普通 2员工券
		DiscountType: 2, //1满减 2折扣
		Discount:     0,
		DiscountFrom: 20,
		Added:        false,
		Ratio:        0.8,
		Used:         false,
		EnableAt:     1666188383,
		ExpiredAt:    1667023583,
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestCouponServer_ListCoupon(t *testing.T) {
	res, err := client.ListCoupon(context.Background(), &pb.ListCouponReq{
		PageSize:     6,
		PageNo:       1,
		Name:         "",
		EnableAt:     0,
		EnableAtOpt:  "",
		ExpiredAt:    0,
		ExpiredAtOpt: "",
		Used:         0,
		Added:        0,
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}
