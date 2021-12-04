package handler

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
	"github.com/i-coder-robot/mic-trainning-lessons/internal"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func CaptchaHandler(c *gin.Context) {

	mobile, ok := c.GetQuery("mobile")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	}

	fileNme := "data.png"
	f, err := os.Create(fileNme)
	if err != nil {
		zap.S().Error("GenCaptcha() 失败")
		return
	}
	defer f.Close()
	var w io.WriterTo
	d := captcha.RandomDigits(captcha.DefaultLen)
	w = captcha.NewImage("", d, captcha.StdWidth, captcha.StdHeight)
	_, err = w.WriteTo(f)
	if err != nil {
		zap.S().Error("GenCaptcha() 失败")
		return
	}
	fmt.Println(d)
	captcha := ""
	for _, item := range d {
		captcha += fmt.Sprintf("%d", item)
	}
	fmt.Println(captcha)
	internal.RedisClient.Set(context.Background(), mobile, captcha, 120*time.Second)
	b64, err := GetBase64(fileNme)
	if err != nil {
		zap.S().Error("GenCaptcha() 失败")
		return
	}
	fmt.Println(b64)
	c.JSON(http.StatusOK, gin.H{
		"captcha": b64,
	})
}

func GetBase64(fileName string) (string, error) {
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", err
	}

	b := make([]byte, 10240)
	base64.StdEncoding.Encode(b, file)
	s := string(b)
	return s, nil
}

func GetBase64_2(fileName string) (string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
		return "", nil
	}
	defer file.Close()

	fileinfo, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return "", nil
	}

	filesize := fileinfo.Size()
	//buffer := make([]byte, filesize)
	buffer := make([]byte, base64.StdEncoding.EncodedLen(int(filesize)))

	bytesread, err := file.Read(buffer)
	if err != nil {
		fmt.Println(err)
		return "", nil
	}

	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Print(err)
	}
	fmt.Println("bytes read: ", bytesread)
	fmt.Println("bytestream to string: ", string(buffer))
	base64.StdEncoding.Encode(buffer, b)
	s := string(b)
	return s, nil
}
