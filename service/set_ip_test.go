package service

import (
	"fmt"
	"testing"
	"yanglu/config"
	"yanglu/service/data"
	"yanglu/service/logger"
)

func TestSetIp(t *testing.T) {
	config.InitEnvConf("../config/env.toml")

	// init log
	logger.InitLogger("", nil)

	// init db
	data.InitMysql()
	ss := NewSetIpService("172.0.0.1")
	err := ss.Ubuntu()
	fmt.Println("err = ", err)
}
