package service

import (
	"testing"
	"yanglu/config"
	"yanglu/service/data"
	"yanglu/service/logger"
)

func TestStartTask(t *testing.T) {
	config.InitEnvConf("../config/env.toml")

	// init log
	logger.InitLogger("", nil)

	// init db
	data.InitMysql()

	ts := NewTaskService()

	task, _ := ts.AddFastTask("47.104.213.134")
	ts.ExecuteTask(task)
}
