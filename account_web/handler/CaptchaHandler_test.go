package handler

import (
	"fmt"
	"testing"
)

func TestGenCaptcha(t *testing.T) {
	//err := GenCaptcha()
	//if err != nil {
	//	panic(err)
	//}
}

func TestGetBase64_2(t *testing.T) {
	fileNme := "data.png"
	s, err := GetBase64_2(fileNme)
	if err != nil {
		panic(err)
	}
	fmt.Println(s)
}
