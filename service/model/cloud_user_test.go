package model

import (
	"fmt"
	"testing"
	"time"
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

	cloudUser := CloudUser{
		Company:    "哈哈哈",
		Phone:      "13112345678",
		Email:      "123@qq.com",
		PassWd:     "12345",
		Authority:  Ints{2},
		CreateTime: time.Now().Unix(),
	}

	err := cloudUser.Create()
	fmt.Println("err = ", err)
	fmt.Println("id = ", cloudUser.Uid)
}
