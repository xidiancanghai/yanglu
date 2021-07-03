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
	ss := NewSetIpService("192.168.1.126")
	// data, err := ss.Gateway4()
	// fmt.Println("err = ", err, " data = ", string(data))

	// is, err := ss.IpHasUsed()
	// fmt.Println("err = ", err, " is = ", is)

	// server, err := ss.GetNameServer()
	// fmt.Println("err = ", err, " server = ", server)

	err := ss.SetIp()
	fmt.Println("err = ", err)
}
