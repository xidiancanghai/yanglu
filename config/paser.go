package config

import (
	"log"
	"path"
	"runtime"

	"github.com/BurntSushi/toml"
)

type EnvType int8

const (
	EnvUnKnown EnvType = 0
	EnvLocal   EnvType = 1
	EnvDevelop EnvType = 2
	EnvOnline  EnvType = 3
)

func (env EnvType) ToString() string {
	switch env {
	case EnvUnKnown:
		return "EnvUnKnown"
	case EnvLocal:
		return "EnvLocal"
	case EnvDevelop:
		return "EnvDevelop"
	case EnvOnline:
		return "EnvOnline"
	default:
		return "EnvNil"
	}
}

type Email struct {
	User   string `toml:"user"`
	Passwd string `toml:"passwd"`
	Addr   string `toml:"addr"`
	Host   string `toml:"host"`
}

type Config struct {
	Environment struct {
		Env     EnvType `toml:"env"`
		LogPath string  `toml:"log_path"`
	} `toml:"environment"`
	Mysql struct {
		Soft struct {
			Dsn string `toml:"dsn"`
		} `toml:"soft"`
		Cloud struct {
			Dsn string `toml:"dsn"`
		} `toml:"cloud"`
	} `toml:"mysql"`
	Http struct {
		Port int `toml:"port"`
	} `toml:"http"`

	AdminUser struct {
		User   string `toml:"user"`
		PassWd string `toml:"pass_wd"`
	} `toml:"amdin_user"`
	Email Email `toml:"email"`

	//Alert struct {
	//	Open bool     `toml:"open"`
	//	At   []string `toml:"at"`
	//} `toml:"alert"`
	//
	//Wechat struct {
	//	MiniProgramAccessToken string `toml:"miniProgramToken"`
	//} `toml:"wechat"`
}

var conf Config

// 仅供单元测试使用
func GetEnvPath() string {
	_, filename, _, _ := runtime.Caller(0)
	return path.Dir(filename)
}

func InitTestConfig() {
	InitEnvConf(GetEnvPath() + "/env.toml")
}

func InitEnvConf(file ...string) {
	configPath := "../config/env.toml"
	if len(file) > 0 {
		configPath = file[0]
	}

	if _, err := toml.DecodeFile(configPath, &conf); err != nil {
		log.Fatalf("failed to decode file, path: %v, err: %v", configPath, err)
	}
	if GetEnv() == EnvUnKnown {
		log.Fatalf("invalied env type: %v", GetEnv())
	}
	log.Printf("succeed to parse env config: %+v", conf)
}

func IsTest() bool {
	if conf.Environment.Env == EnvLocal || conf.Environment.Env == EnvDevelop {
		return true
	}
	return false
}

func IsLocal() bool {
	if conf.Environment.Env == EnvLocal {
		return true
	}
	return false
}

func IsOnline() bool {
	if conf.Environment.Env == EnvOnline {
		return true
	}
	return false
}

// local, develop, online
func GetEnv() EnvType {
	return conf.Environment.Env
}

func GetLogPath() string {
	return conf.Environment.LogPath
}

func GetMysqlApi() string {
	if IsCloud() {
		return conf.Mysql.Cloud.Dsn
	}
	return conf.Mysql.Soft.Dsn
}

func GetHttpPort() int {
	return conf.Http.Port
}

func GetAdminInfo() (string, string) {
	return conf.AdminUser.User, conf.AdminUser.PassWd
}

func GetEmailConf() *Email {
	return &conf.Email
}
