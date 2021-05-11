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

	user := User{
		Name:       "xionger",
		Passwd:     "1234567",
		Authority:  0,
		Department: "安全",
	}
	err := user.Create()
	fmt.Println("err = ", err)
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
