package biz

import (
	"crypto/md5"
	"fmt"
	"github.com/anaskhan96/go-password-encoder"
	"testing"
)

func TestGetMd5(t *testing.T) {
	//s:="happy"
	//fmt.Println(GetMd5(s))
	//
	//s="happy1"
	//fmt.Println(GetMd5(s))
	//
	//s="happy2"
	//fmt.Println(GetMd5(s))

	//s=fmt.Sprintf("%s%d",s,time.Now().Unix())
	//fmt.Println(GetMd5(s))
	//time.Sleep(2*time.Second)
	//s=fmt.Sprintf("%s%d",s,time.Now().Unix())
	//fmt.Println(GetMd5(s))

	//salt, encodedPwd := password.Encode("happy", nil)
	//fmt.Println(salt)
	//fmt.Println(encodedPwd)
	//
	//check := password.Verify("happy1", salt, encodedPwd, nil)
	//fmt.Println(check)

	options := password.Options{
		SaltLen:      16,
		Iterations:   100,
		KeyLen:       32,
		HashFunction: md5.New,
	}
	salt, hashedPwd := password.Encode("happy", &options)
	fmt.Println(salt)
	fmt.Println(hashedPwd)

	check := password.Verify("happy", salt, hashedPwd, &options)
	fmt.Println(check)
}
