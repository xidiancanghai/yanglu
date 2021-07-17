package service

import (
	"fmt"
	"testing"
	"time"
	"yanglu/config"
	"yanglu/service/data"
	"yanglu/service/logger"
)

func TestStartTask(t *testing.T) {

	config.InitLicenseConfig()
	config.InitEnvConf("../config/env.toml")

	// init log
	logger.InitLogger("", nil)

	// init db
	data.InitMysql()

	ts := NewTaskService()

	task, _ := ts.AddFastTask("47.104.213.134")
	err := ts.ExecuteTask(task)
	fmt.Println("err = ", err)
	time.Sleep(120 * time.Second)
}
