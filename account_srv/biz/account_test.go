package biz

import (
	"context"
	"fmt"
	"github.com/i-coder-robot/mic-trainning-lessons/account_srv/internal"
	"github.com/i-coder-robot/mic-trainning-lessons/account_srv/proto/pb"
	"testing"
)

func init() {
	internal.InitDB()
}

func TestAccountServer_AddAccount(t *testing.T) {
	accountServer := AccountServer{}
	for i := 0; i < 5; i++ {
		s := fmt.Sprintf("1300000000%d", i)
		res, err := accountServer.AddAccount(context.Background(), &pb.AddAccountRequest{
			Mobile:   s,
			Password: s,
			NickName: s,
			Gender:   "male",
		})
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(res.Id)
	}
}

func TestAccountServer_GetAccountList(t *testing.T) {
	accountServer := AccountServer{}
	res, err := accountServer.GetAccountList(context.Background(), &pb.PagingRequest{
		PageNo:   1,
		PageSize: 3,
	})
	if err != nil {
		fmt.Println(err)
	}
	for _, account := range res.AccountList {
		fmt.Println(account.Id)
	}

	res, err = accountServer.GetAccountList(context.Background(), &pb.PagingRequest{
		PageNo:   2,
		PageSize: 3,
	})
	if err != nil {
		fmt.Println(err)
	}
	for _, account := range res.AccountList {
		fmt.Println(account.Id)
	}
}

func TestAccountServer_GetAccountByMobile(t *testing.T) {
	mobile := "13000000000"
	accountServer := AccountServer{}
	res, err := accountServer.GetAccountByMobile(context.Background(), &pb.MobileRequest{Mobile: mobile})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res.Id)
}

func TestAccountServer_GetAccountById(t *testing.T) {
	id := 3
	accountServer := AccountServer{}
	res, err := accountServer.GetAccountById(context.Background(), &pb.IdRequest{Id: uint32(id)})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res.Mobile)
}

func TestAccountServer_UpdateAccount(t *testing.T) {
	accountServer := AccountServer{}
	req := pb.UpdateAccountRequest{
		Id:       1,
		Mobile:   "13000000100",
		Password: "13000000100",
		NickName: "13000000100",
		Gender:   "female",
		Role:     2,
	}
	res, err := accountServer.UpdateAccount(context.Background(), &req)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res.Result)
}

func TestAccountServer_CheckPassword(t *testing.T) {
	accountServer := AccountServer{}
	res, err := accountServer.CheckPassword(context.Background(), &pb.CheckPasswordRequest{
		Password:       "13000000004",
		HashedPassword: "e06e418810aa3c411a7fbe623ea4e8d9338740fcd5eb3996c7163f189c999b79",
		AccountId:      5,
	})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res.Result)
}
