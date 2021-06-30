package model

import (
	"fmt"
	"testing"
	"yanglu/config"
	"yanglu/service/data"
	"yanglu/service/logger"
)

func TestCloudUser(t *testing.T) {
	config.InitEnvConf("../../config/env.toml")

	// init log
	logger.InitLogger("", nil)

	// init db
	data.InitMysql()

	// cloudUser := &CloudUser{
	// 	Company:    "哈哈哈",
	// 	Phone:      "13112345678",
	// 	Email:      "123@qq.com",
	// 	Passwd:     "12345",
	// 	Authority:  Ints{2},
	// 	CreateTime: time.Now().Unix(),
	// }

	// err := cloudUser.Create()
	// fmt.Println("err = ", err)
	// fmt.Println("id = ", cloudUser.Uid)
	u, err := NewCloudUser().GetUser(map[string]interface{}{
		"phone": "13112345678",
	})
	fmt.Println("u = ", *u, " err = ", err)
}
