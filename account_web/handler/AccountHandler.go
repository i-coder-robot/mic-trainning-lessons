package handler

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/i-coder-robot/mic-trainning-lessons/account_srv/proto/pb"
	"github.com/i-coder-robot/mic-trainning-lessons/account_web/res"
	"github.com/i-coder-robot/mic-trainning-lessons/custom_error"
	"github.com/i-coder-robot/mic-trainning-lessons/log"
	"google.golang.org/grpc"
	"net/http"
	"strconv"
)

func HandleError(err error) string {
	if err != nil {
		switch err.Error() {
		case custom_error.AccountExists:
			return custom_error.AccountExists
		case custom_error.AccountNotFound:
			return custom_error.AccountNotFound
		case custom_error.SaltError:
			return custom_error.SaltError
		default:
			return custom_error.InternalError
		}
	}
	return ""
}

func AccountListHandler(c *gin.Context) {
	pageNoStr := c.DefaultQuery("pageNo", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "3")
	conn, err := grpc.Dial("127.0.0.1:9095", grpc.WithInsecure())
	if err != nil {
		s := fmt.Sprintf("AccountListHandler-GRPC拨号失败:%s", err.Error())
		log.Logger.Info(s)
		e := HandleError(err)
		c.JSON(http.StatusOK, gin.H{
			"msg": e,
		})
		return
	}
	pageNo, _ := strconv.ParseInt(pageNoStr, 10, 32)
	pageSize, _ := strconv.ParseInt(pageSizeStr, 10, 32)
	client := pb.NewAccountServiceClient(conn)
	r, err := client.GetAccountList(context.Background(), &pb.PagingRequest{
		PageNo:   int32(pageNo),
		PageSize: int32(pageSize),
	})
	if err != nil {
		s := fmt.Sprintf("GetAccountList调用失败:%s", err.Error())
		log.Logger.Info(s)
		e := HandleError(err)
		c.JSON(http.StatusOK, gin.H{
			"msg": e,
		})
		return
	}
	var resList []res.Account4Res
	for _, item := range r.AccountList {
		resList = append(resList, pb2res(item))
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":   "",
		"total": r.Total,
		"data":  resList,
	})
}

func pb2res(accountRes *pb.AccountRes) res.Account4Res {
	return res.Account4Res{
		Mobile:   accountRes.Mobile,
		NickName: accountRes.Nickname,
		Gender:   accountRes.Gender,
	}
}
