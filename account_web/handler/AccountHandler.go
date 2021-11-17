package handler

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/i-coder-robot/mic-trainning-lessons/account_srv/proto/pb"
	"github.com/i-coder-robot/mic-trainning-lessons/account_web/req"
	"github.com/i-coder-robot/mic-trainning-lessons/account_web/res"
	"github.com/i-coder-robot/mic-trainning-lessons/custom_error"
	"github.com/i-coder-robot/mic-trainning-lessons/jwt_op"
	"github.com/i-coder-robot/mic-trainning-lessons/log"
	"google.golang.org/grpc"
	"net/http"
	"strconv"
	"time"
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

func LoginByPasswordHandler(c *gin.Context) {
	var loginByPassword req.LoginByPassword
	err := c.ShouldBindJSON(&loginByPassword)
	if err != nil {
		log.Logger.Error("LoginByPassword出错：" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"msg": "解析参数错误",
		})
		return
	}
	//TODO 校验手机号码格式
	//loginByPassword.Mobile不匹配正则表达式，就报错
	conn, err := grpc.Dial("127.0.0.1:9095", grpc.WithInsecure())
	if err != nil {
		log.Logger.Error("LoginByPasswordHandler 拨号出错：" + err.Error())
		e := HandleError(err)
		c.JSON(http.StatusOK, gin.H{
			"msg": e,
		})
		return
	}
	client := pb.NewAccountServiceClient(conn)
	r, err := client.GetAccountByMobile(context.Background(), &pb.MobileRequest{Mobile: loginByPassword.Mobile})
	if err != nil {
		log.Logger.Error("GRPC GetAccountByMobile 出错：" + err.Error())
		e := HandleError(err)
		c.JSON(http.StatusOK, gin.H{
			"msg": e,
		})
		return
	}
	cheRes, err := client.CheckPassword(context.Background(), &pb.CheckPasswordRequest{
		Password:       loginByPassword.Password,
		HashedPassword: r.Password,
		AccountId:      uint32(r.Id),
	})
	if err != nil {
		log.Logger.Error("GRPC CheckPassword 出错：" + err.Error())
		e := HandleError(err)
		c.JSON(http.StatusOK, gin.H{
			"msg": e,
		})
		return
	}
	checkResult := "登录失败"
	if cheRes.Result {
		checkResult = "登录成功"
		j := jwt_op.NewJWT()
		now := time.Now()
		claims := jwt_op.CustomClaims{
			StandardClaims: jwt.StandardClaims{
				NotBefore: now.Unix(),
				ExpiresAt: now.Add(time.Hour * 24 * 30).Unix(),
			},
			ID:          r.Id,
			NickName:    r.Nickname,
			AuthorityId: int32(r.Role),
		}
		token, err := j.GenerateJWT(claims)
		if err != nil {
			log.Logger.Error("GRPC GenerateJWT 出错：" + err.Error())
			e := HandleError(err)
			c.JSON(http.StatusOK, gin.H{
				"msg": e,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"msg":    "",
			"result": checkResult,
			"token":  token,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":    "",
		"result": checkResult,
		"token":  "",
	})
}

//eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9
//.eyJleHAiOjE2Mzk0NjU4MzYsIm5iZiI6MTYzNjg3MzgzNiwiSUQiOjUsIk5pY2tOYW1lIjoiMTMwMDAwMDAwMDQiLCJBdXRob3JpdHlJZCI6MX0
//.47YspOTF5kGO84KMN56ksJzC6sAcMCtqp13D00X6ZBI
