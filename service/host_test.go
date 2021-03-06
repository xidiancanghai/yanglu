package service

import (
	"fmt"
	"testing"
	"yanglu/config"
	"yanglu/service/data"
	"yanglu/service/logger"
	"yanglu/service/model"
)

func TestHostInfo(t *testing.T) {
	config.InitEnvConf("../config/env.toml")

	// init log
	logger.InitLogger("", nil)

	// init db
	data.InitMysql()

	hs := NewHostInfoService()
	//_, err := hs.Add("112.125.25.235", 22, "testlsm", "testlsm123", "管理")
	_, err := hs.Add("192.168.1.1", 110, "admin", "admin", "", "hello")
	fmt.Println("err = ", err)
}

func TestCheckPass(t *testing.T) {
	hostInfo := model.HostInfo{
		Ip:        "112.125.25.235",
		Port:      22,
		SshUser:   "testlsm",
		SshPasswd: "testlsm123",
	}
	hs := NewHostInfoService()
	hs.CheckPass(&hostInfo)
}

func TestBatchCheck(t *testing.T) {
	hostInfo := model.HostInfo{
		Ip:        "112.125.25.235",
		Port:      22,
		SshUser:   "testlsm",
		SshPasswd: "testlsm123",
	}
	list := []*model.HostInfo{}
	list = append(list, &hostInfo)
	list = append(list, &model.HostInfo{
		Ip:        "112.125.25.235",
		Port:      22,
		SshUser:   "testlsm",
		SshPasswd: "testlsm1234",
	})
	hs := NewHostInfoService()
	err := hs.BatchCheck(list)
	fmt.Println("err = ", err)
}

func TestDepartment(t *testing.T) {
	config.InitEnvConf("../config/env.toml")

	// init log
	logger.InitLogger("", nil)

	// init db
	data.InitMysql()

	hs := NewHostInfoService()
	err := hs.UpdateDepartment("112.125.25.235", "测试部门")
	fmt.Println("err = ", err)
}

func TestGetSystemInfo(t *testing.T) {
	config.InitEnvConf("../config/env.toml")

	// init log
	logger.InitLogger("", nil)

	// init db
	data.InitMysql()

	hs := NewHostInfoService()

	hostInfo, _ := model.NewHostInfo().GetHostInfoByIp("112.125.25.235")

	res, err := hs.GetSystemInfo(hostInfo)
	fmt.Println("err = ", err, " res = ", res)
}

func TestUpdateSystemInfo(t *testing.T) {
	config.InitEnvConf("../config/env.toml")

	// init log
	logger.InitLogger("", nil)

	// init db
	data.InitMysql()

	hs := NewHostInfoService()

	hostInfo, _ := model.NewHostInfo().GetHostInfoByIp("47.104.213.134")

	err := hs.Prepare([]*model.HostInfo{hostInfo})
	fmt.Println("err = ", err)

	//fmt.Println(hs.CheckDir(hostInfo, "/home/ftp1"))
}

func TestCpFile(t *testing.T) {
	config.InitEnvConf("../config/env.toml")

	// init log
	logger.InitLogger("", nil)

	// init db
	data.InitMysql()

	hs := NewHostInfoService()

	hostInfo, _ := model.NewHostInfo().GetHostInfoByIp("112.125.25.235")

	//err := hs.CpFile(hostInfo, "trivy_0.16.0_Linux-64bit.tar.gz", "/var/trivy")

	err := hs.CpFileBySftp(hostInfo, "trivy-offline.db.tgz", "/root/.cache/trivy/db")

	fmt.Println(" err = ", err)
}

func TestCmd(t *testing.T) {

	config.InitLicenseConfig()

	config.InitEnvConf("../config/env.toml")

	// init log
	logger.InitLogger("", nil)

	// init db
	data.InitMysql()

	hs := NewHostInfoService()

	_, err := hs.Add("47.104.213.134", 22, "testyly", "testyly@123", "测试", "测试")
	fmt.Println("err = ", err)
}

func TestPrepare(t *testing.T) {

	config.InitLicenseConfig()

	config.InitEnvConf("../config/env.toml")

	// init log
	logger.InitLogger("", nil)

	// init db
	data.InitMysql()

	hs := NewHostInfoService()

	hostInfo, _ := model.NewHostInfo().GetHostInfoByIp("47.104.213.134")

	hs.Prepare([]*model.HostInfo{hostInfo})
}

func TestGetHostName(t *testing.T) {

	config.InitLicenseConfig()

	config.InitEnvConf("../config/env.toml")

	// init log
	logger.InitLogger("", nil)

	// init db
	data.InitMysql()

	hs := NewHostInfoService()

	hostInfo, _ := model.NewHostInfo().GetHostInfoByIp("47.104.213.134")
	fmt.Println(hs.GetHostName(hostInfo))

	hostInfo1, _ := model.NewHostInfo().GetHostInfoByIp("112.125.25.235")
	fmt.Println(hs.GetHostName(hostInfo1))
}
