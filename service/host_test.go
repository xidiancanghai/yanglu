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

	hostInfo, _ := model.NewHostInfo().GetHostInfoByIp("112.125.25.235")

	err := hs.Prepare([]*model.HostInfo{hostInfo})
	fmt.Println("err = ", err)
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
	config.InitEnvConf("../config/env.toml")

	// init log
	logger.InitLogger("", nil)

	// init db
	data.InitMysql()

	hs := NewHostInfoService()

	hostInfo, _ := model.NewHostInfo().GetHostInfoByIp("112.125.25.235")

	//hs.CpFile(hostInfo, "cmd.sh")

	//res, err := hs.Cmd(hostInfo, "bash /var/cmd.sh")
	//_, err := hs.Cmd(hostInfo, "cd /var/trivy;tar -xzvf trivy_0.16.0_Linux-64bit.tar.gz")

	res, err := hs.Cmd(hostInfo, "cd /var/trivy; rm -rf results.json;./trivy fs -f json -o results.json /")

	fmt.Println(" res = ", res, " err ", err)

	res, err = hs.GetResult(hostInfo)
	fmt.Println(" res = ", res, " err ", err)
}
