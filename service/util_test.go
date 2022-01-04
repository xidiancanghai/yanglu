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

func TestBuildPdf(t *testing.T) {
	// us := NewUtilService()

	// list := []map[string]interface{}{}
	// list = append(list, map[string]interface{}{
	// 	"a": 1,
	// 	"b": 2,
	// })
	// us.buildPdf(list)

}
