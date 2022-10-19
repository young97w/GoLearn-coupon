package biz

import (
	"context"
	"coupon/coupon_srv/model"
	"coupon/custom_error"
	"coupon/internal"
	"coupon/log"
	"coupon/proto/pb"
	"errors"
	"sync"
)

var m sync.Mutex

type CouponServer struct {
}

func (c CouponServer) AddCoffee(ctx context.Context, req *pb.CoffeeItem) (*pb.CoffeeItem, error) {
	var coffee model.Coffee
	convertCoffeePb2Model(req, &coffee)
	r := internal.DB.Save(&coffee)
	if r.RowsAffected == 0 {
		log.Logger.Error(custom_error.AddCoffeeFailed)
		return nil, errors.New(custom_error.AddCoffeeFailed)
	}
	return req, nil
}

func (c CouponServer) DeleteCoffee(ctx context.Context, req *pb.DeleteCoffeeReq) (*pb.DeleteCoffeeRes, error) {
	r := internal.DB.Delete(&model.Coffee{}, req.Id)
	if r.RowsAffected == 0 {
		log.Logger.Error(custom_error.DeleteCoffeeFailed)
		return nil, errors.New(custom_error.DeleteCoffeeFailed)
	}
	return &pb.DeleteCoffeeRes{Result: true}, nil
}

func (c CouponServer) UpdateCoffee(ctx context.Context, req *pb.CoffeeItem) (*pb.CoffeeItem, error) {
	m.Lock()
	defer m.Unlock()
	var coffee model.Coffee
	r := internal.DB.First(&coffee, req.Id)
	if r.RowsAffected == 0 {
		log.Logger.Error(custom_error.CannotFindCoffee)
		return nil, errors.New(custom_error.CannotFindCoffee)
	}
	convertCoffeePb2Model(req, &coffee)
	r = internal.DB.Save(&coffee)
	if r.RowsAffected == 0 {
		log.Logger.Error(custom_error.UpdateCoffeeFailed)
		return nil, errors.New(custom_error.UpdateCoffeeFailed)
	}
	return req, nil
}

func (c CouponServer) SellCoffee(ctx context.Context, req *pb.SellCoffeeReq) (*pb.CoffeeItem, error) {
	m.Lock()
	defer m.Unlock()
	var coffee model.Coffee
	r := internal.DB.First(&coffee, req.Id)
	if r.RowsAffected == 0 {
		log.Logger.Error(custom_error.CannotFindCoffee)
		return nil, errors.New(custom_error.CannotFindCoffee)
	}
	tx := internal.DB.Begin()
	coffee.SoldNum += req.SoldNum
	r = tx.Save(coffee)
	if r.RowsAffected == 0 {
		log.Logger.Error(custom_error.UpdateCoffeeFailed)
		tx.Rollback()
		return nil, errors.New(custom_error.UpdateCoffeeFailed)
	}
	tx.Commit()
	return ConvertCoffeeModel2Pb(&coffee), nil
}

func (c CouponServer) ListCoffee(ctx context.Context, req *pb.ListCoffeeReq) (*pb.CoffeeListRes, error) {
	var coffees []model.Coffee
	var coffeeList []*pb.CoffeeItem
	var res pb.CoffeeListRes
	r := internal.DB.Model(&model.Coffee{}).Scopes(internal.Paginate(int(req.PageSize), int(req.PageNo))).Find(&coffees)
	if r.RowsAffected == 0 {
		log.Logger.Error(custom_error.CannotFindCoffee)
		return nil, errors.New(custom_error.CannotFindCoffee)
	}
	for _, coffee := range coffees {
		item := ConvertCoffeeModel2Pb(&coffee)
		coffeeList = append(coffeeList, item)
	}
	res.Total = int32(r.RowsAffected)
	res.CoffeeList = coffeeList
	return &res, nil
}

func (c CouponServer) GetCoffee(ctx context.Context, req *pb.CoffeeItem) (*pb.CoffeeItem, error) {
	var coffee model.Coffee
	var res *pb.CoffeeItem
	r := internal.DB.First(&coffee, req.Id)
	if r.RowsAffected == 0 {
		log.Logger.Error(custom_error.CannotFindCoffee)
		return nil, errors.New(custom_error.CannotFindCoffee)
	}
	res = ConvertCoffeeModel2Pb(&coffee)
	return res, nil
}

func convertCoffeePb2Model(coffeeReq *pb.CoffeeItem, coffee *model.Coffee) {
	if coffeeReq.Id != 0 {
		coffee.ID = uint(coffeeReq.Id)
	}
	if len(coffeeReq.Name) != 0 {
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
	if len(coffeeReq.Sku) != 0 {
		coffee.Sku = coffeeReq.Sku
	}
	if len(coffeeReq.Description) != 0 {
		coffee.Description = coffeeReq.Description
	}
	if len(coffeeReq.Image) != 0 {
		coffee.Image = coffeeReq.Image
	}
}

func ConvertCoffeeModel2Pb(coffee *model.Coffee) *pb.CoffeeItem {
	item := pb.CoffeeItem{
		Id:          int32(coffee.ID),
		Name:        coffee.Name,
		Price:       coffee.Price,
		RealPrice:   coffee.RealPrice,
		SoldNum:     coffee.SoldNum,
		Sku:         coffee.Sku,
		Description: coffee.Description,
		Image:       coffee.Image,
	}

	return &item
}
