package model

import (
	"fmt"
	"testing"
	"yanglu/config"
	"yanglu/service/data"
	"yanglu/service/logger"
)

func TestTask(t *testing.T) {
	config.InitEnvConf("../../config/env.toml")

	// init log
	logger.InitLogger("", nil)

	// init db
	data.InitMysql()

	// task := &Task{
	// 	Ip:        "112.125.25.235",
	// 	IsRepeate: 1,
	// 	ExecuTime: "9h30min",
	// }
	// task.Create()
	rt, _ := NewTask().GetTask(1)
	fmt.Println("task = ", *rt)
	rt.Updates(map[string]interface{}{
		"status": 1,
	})
}

func TestGetMaxBatchId(t *testing.T) {
	config.InitEnvConf("../../config/env.toml")

	// init log
	logger.InitLogger("", nil)

	// init db
	data.InitMysql()

	// task := &Task{
	// 	Ip:        "112.125.25.235",
	// 	IsRepeate: 1,
	// 	ExecuTime: "9h30min",
	// }
	// task.Create()
	// rt, _ := NewTask().GetMaxBatchId()
	// fmt.Println("task = ", rt)

}

func TestGetSmartTask(t *testing.T) {
	config.InitEnvConf("../../config/env.toml")

	// init log
	logger.InitLogger("", nil)

	// init db
	data.InitMysql()

	list, _ := NewTask().GetSmartTask()
	for _, v := range list {
		fmt.Println("item = ", *v)
	}
}
