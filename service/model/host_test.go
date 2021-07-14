package model

import (
	"fmt"
	"testing"
	"yanglu/config"
	"yanglu/service/data"
	"yanglu/service/logger"
)

func TestHostInfo(t *testing.T) {
	config.InitEnvConf("../../config/env.toml")

	// init log
	logger.InitLogger("", nil)

	// init db
	data.InitMysql()

	hostInfo := HostInfo{
		Ip:        "127.0.0.1",
		Port:      22,
		SshUser:   "canghai",
		SshPasswd: "canghai123",
	}
	err := hostInfo.Create()
	fmt.Println("err = ", err)
}

func TestGetHostInfoById(t *testing.T) {
	config.InitEnvConf("../../config/env.toml")

	// init log
	logger.InitLogger("", nil)

	// init db
	data.InitMysql()

	host, err := NewHostInfo().GetHostsByIps([]string{"112.125.25.235"})
	fmt.Println("err = ", err, " host = ", host)
}
