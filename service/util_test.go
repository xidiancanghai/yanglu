package service

import (
	"fmt"
	"testing"
)

func TestGetCaptCha(t *testing.T) {
	us := NewUtilService()
	buf, err := us.GetCaptcha(us.GetCaptchaId())
	fmt.Println("buf = ", buf, " err = ", err)
}

// func TestUploadImage(t *testing.T) {
// 	us := NewUtilService()
// 	us.UploadImage("xiaolongnv.jpg")
// }
