package biz

import (
	"account/account_srv/model"
	"account/custom_error"
	"account/internal"
	"account/log"
	"account/password"
	"account/proto/pb"
	"context"
	"errors"
	"sync"
)

var m sync.Mutex

type AccountServer struct {
}

func (a AccountServer) GetAccountList(ctx context.Context, req *pb.ListAccountReq) (*pb.AccountListRes, error) {
	var accounts []model.Account
	var res pb.AccountListRes
	var accountList []*pb.AccountRes

	r := internal.DB.Scopes(internal.Paginate(int(req.PageSize), int(req.PageNo))).Find(&accounts)
	if r.RowsAffected < 1 {
		log.Logger.Error(custom_error.AccountGetFailed)
		return nil, errors.New(custom_error.AccountGetFailed)
	}
	for _, item := range accounts {
		accountPb := convertAccountModel2Pb(&item)
		accountList = append(accountList, &accountPb)
	}
	res.Total = int32(r.RowsAffected)
	res.AccountList = accountList
	return &res, nil
}

func (a AccountServer) GetAccountByMobile(ctx context.Context, req *pb.MobileAccountReq) (*pb.AccountRes, error) {
	var account model.Account
	if req.Mobile == "" {
		log.Logger.Error(custom_error.ParameterIncorrect)
		return nil, errors.New(custom_error.ParameterIncorrect)
	}
	r := internal.DB.Where("mobile=?", req.Mobile).First(&account)
	if r.RowsAffected < 1 {
		log.Logger.Error(custom_error.AccountNotFind)
		return nil, errors.New(custom_error.AccountNotFind)
	}
	res := convertAccountModel2Pb(&account)
	return &res, nil
}

func (a AccountServer) GetAccountById(ctx context.Context, req *pb.IdAccountReq) (*pb.AccountRes, error) {
	var account model.Account
	if req.Id == 0 {
		log.Logger.Error(custom_error.ParameterIncorrect)
		return nil, errors.New(custom_error.ParameterIncorrect)
	}
	r := internal.DB.First(&account, req.Id)
	if r.RowsAffected < 1 {
		log.Logger.Error(custom_error.AccountNotFind)
		return nil, errors.New(custom_error.AccountNotFind)
	}
	res := convertAccountModel2Pb(&account)
	return &res, nil
}

func (a AccountServer) AddAccount(ctx context.Context, req *pb.AddAccountReq) (*pb.AccountRes, error) {
	var account model.Account
	r := internal.DB.Model(model.Account{}).Where("mobile=?", req.Mobile).First(&account)
	if r.RowsAffected > 0 {
		log.Logger.Error(custom_error.AccountAlreadyExists)
		return nil, errors.New(custom_error.AccountAlreadyExists)
	}
	salt, hsdPwd := password.GenerateHashedPwd(req.Password)
	account.Salt = salt
	account.Password = hsdPwd
	account.Role = int(req.Role)
	account.Mobile = req.Mobile
	account.NickName = req.Nickname
	account.Gender = req.Gender
	account.IsEmployee = req.IsEmployee
	r = internal.DB.Save(&account)
	if r.RowsAffected == 0 {
		log.Logger.Error(custom_error.AccountCreateFailed)
		return nil, errors.New(custom_error.AccountCreateFailed)
	}
	pb := convertAccountModel2Pb(&account)
	return &pb, nil
}

func (a AccountServer) UpdateAccount(ctx context.Context, req *pb.UpdateAccountReq) (*pb.UpdateAccountRes, error) {
	var account model.Account
	m.Lock()
	defer m.Unlock()
	r := internal.DB.Model(model.Account{}).First(&account, req.Id)
	if r.RowsAffected == 0 {
		log.Logger.Error(custom_error.AccountNotFind)
		return nil, errors.New(custom_error.AccountNotFind)
	}
	if req.Role != 0 {
		account.Role = int(req.Role)
	}
	if req.Nickname != "" {
		account.NickName = req.Nickname
	}
	if req.Mobile != "" {
		account.Mobile = req.Mobile
	}
	if req.Gender != "" {
		account.Gender = req.Gender
	}
	r = internal.DB.Save(&account)
	if r.RowsAffected == 0 {
		log.Logger.Error(custom_error.InternalError)
		return nil, errors.New(custom_error.InternalError)
	}
	return &pb.UpdateAccountRes{Result: true}, nil
}

func (a AccountServer) CheckPassword(ctx context.Context, req *pb.CheckPasswordReq) (*pb.CheckPasswordRes, error) {
	var account model.Account
	r := internal.DB.Model(model.Account{}).Where("mobile=?", req.Mobile).First(&account)
	if r.RowsAffected == 0 {
		log.Logger.Error(custom_error.AccountNotFind)
		return nil, errors.New(custom_error.AccountNotFind)
	}
	check := password.CheckPwd(req.Password, account.Salt, account.Password)
	return &pb.CheckPasswordRes{Result: check}, nil
}

func convertAccountModel2Pb(account *model.Account) pb.AccountRes {
	return pb.AccountRes{
		Id:         int32(account.ID),
		Mobile:     account.Mobile,
		Nickname:   account.NickName,
		Gender:     account.Gender,
		IsEmployee: account.IsEmployee,
		Role:       uint32(account.Role),
	}
}
