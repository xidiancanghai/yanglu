package model

import (
	"fmt"
	"testing"
	"yanglu/config"
	"yanglu/service/data"
	"yanglu/service/logger"
)

func TestUser(t *testing.T) {
	config.InitEnvConf("../../config/env.toml")

	// init log
	logger.InitLogger("", nil)

	// init db
	data.InitMysql()

	// user := &User{
	// 	Name:       "xionger",
	// 	Passwd:     "1234567",
	// 	Authority:  []int{1},
	// 	Department: "安全",
	// }
	user := new(User)
	user.Name = "xionger"
	user.Passwd = "1234567"
	err := user.Create()
	fmt.Println("err = ", err, " uid = ", user.Uid)
}

func TestSelect(t *testing.T) {
	config.InitEnvConf("../../config/env.toml")

	// init log
	logger.InitLogger("", nil)

	// init db
	data.InitMysql()

	user, err := NewUser().GetUserByName("xionger")

	fmt.Println("err = ", err)
	fmt.Println("user = ", *user)
}
