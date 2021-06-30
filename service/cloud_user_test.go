package service

import (
	"fmt"
	"testing"
	"yanglu/config"
	"yanglu/service/data"
	"yanglu/service/logger"
)

func TestSendEmail(t *testing.T) {

	config.InitEnvConf("../config/env.toml")

	// init log
	logger.InitLogger("", nil)

	// init db
	data.InitMysql()

	cs := NewEmptyCloudUserService()
	_, err := cs.Register("哈哈", "13152015823", "chmy2272120002@com", "123456")
	if err != nil {
		fmt.Println("err = ", err)
	}
}
