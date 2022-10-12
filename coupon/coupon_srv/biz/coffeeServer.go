package biz

import (
	"context"
	"coupon/coupon_srv/model"
	"coupon/custom_error"
	"coupon/internal"
	"coupon/log"
	"coupon/proto/pb"
	"errors"
)

type CoffeeServer struct {
}

func (c CoffeeServer) AddCoffee(ctx context.Context, req *pb.CoffeeItem) (*pb.CoffeeItem, error) {
	var coffee model.Coffee
	convertCoffeePb2Model(req, &coffee)
	r := internal.DB.Save(&coffee)
	if r.RowsAffected == 0 {
		log.Logger.Error(custom_error.AddCoffeeFailed)
		return nil, errors.New(custom_error.AddCoffeeFailed)
	}
	return req, nil
}

func (c CoffeeServer) DeleteCoffee(ctx context.Context, req *pb.DeleteCoffeeReq) (*pb.DeleteCoffeeRes, error) {
	r := internal.DB.Delete(&model.Coffee{}, req.Id)
	if r.RowsAffected == 0 {
		log.Logger.Error(custom_error.DeleteCoffeeFailed)
		return nil, errors.New(custom_error.DeleteCoffeeFailed)
	}
	return &pb.DeleteCoffeeRes{Result: true}, nil
}

func (c CoffeeServer) UpdateCoffee(ctx context.Context, req *pb.CoffeeItem) (*pb.CoffeeItem, error) {
	//TODO implement me
	panic("implement me")
}

func (c CoffeeServer) ListCoffee(ctx context.Context, req *pb.ListCoffeeReq) (*pb.CoffeeListRes, error) {
	//TODO implement me
	panic("implement me")
}

func (c CoffeeServer) GetCoffee(ctx context.Context, req *pb.CoffeeItem) (*pb.CoffeeItem, error) {
	//TODO implement me
	panic("implement me")
}

func convertCoffeePb2Model(coffeeReq *pb.CoffeeItem, coffee *model.Coffee) {
	if coffeeReq.Id != 0 {
		coffee.ID = uint(coffeeReq.Id)
	}
	if coffeeReq.Name != "" || len(coffeeReq.Name) == 0 {
		coffee.Name = coffeeReq.Name
	}
	if coffeeReq.Price != 0 {
		coffee.Price = coffeeReq.Price
	}
	if coffeeReq.RealPrice != 0 {
		coffee.RealPrice = coffeeReq.RealPrice
	}
	if coffeeReq.SoldNum != 0 {
		coffee.SoldNum = coffeeReq.SoldNum
	}
	if coffeeReq.Sku != "" || len(coffeeReq.Sku) == 0 {
		coffee.Sku = coffeeReq.Sku
	}
	if coffeeReq.Description != "" || len(coffeeReq.Description) == 0 {
		coffee.Description = coffeeReq.Description
	}
	if coffeeReq.Image != "" || len(coffeeReq.Image) == 0 {
		coffee.Image = coffeeReq.Image
	}
}
