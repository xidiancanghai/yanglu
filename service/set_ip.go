package service

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

type SetIpService struct {
	ip string
}

func NewSetIpService(ip string) *SetIpService {
	return &SetIpService{ip: ip}
}

func (ss *SetIpService) SetIp() error {
	cmd := exec.Command("cat", "/etc/issue")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	data, err := cmd.Output()
	data = []byte(`Ubuntu 20.04.2 LTS \n \l`)
	err = nil

	if err != nil {
		logrus.Error("SetIp err ", err)
		return err
	}
	systemOs := strings.ToUpper(string(data))
	isUbuntu := strings.Contains(systemOs, strings.ToUpper("Ubuntu"))
	isCentOS := strings.Contains(systemOs, strings.ToUpper("CentOS"))
	if !isUbuntu && !isCentOS {
		return fmt.Errorf("暂不支持该系统%s设置ip", string(data))
	}
	reg := regexp.MustCompile(`[0-9]+\.[0-9]+`)
	list := reg.FindAllString(systemOs, -1)
	if len(list) == 0 {
		return errors.New("获取系统版本号错误")
	}
	versionNumber, err := strconv.ParseFloat(list[0], 32)
	if err != nil {
		return errors.New("获取系统版本号错误")
	}
	if (isUbuntu && versionNumber < 19.0) || (isCentOS && versionNumber < 7.0) {
		return errors.New("Ubuntu系统只支持19.0以上版本或centos系统支持7.0以上")
	}
	if isUbuntu {
		//ss.Ubuntu()
	}
	return nil
}

// func (ss *SetIpService) Ubuntu() error {
// 	path := "/Users/weipai-liuxiang/50-cloud-init.yaml"
// 	file, err := os.Open(path)
// 	if err != nil {
// 		logrus.Error("Ubuntu err ", err)
// 		return err
// 	}
// 	defer file.Close()
// 	data, err := ioutil.ReadAll(file)
// 	if err != nil {
// 		logrus.Error("Ubuntu err ", err)
// 		return err
// 	}
// 	conf := struct {
// 		Network struct {
// 			Version   string `yaml:"version"`
// 			Ethernets struct {
// 				Eth0 struct {
// 					Dhcp4 bool `yaml:"dhcp4"`
// 					Match struct {
// 						Macaddress string `yaml:"macaddress"`
// 					} `yaml:"match"`
// 					SetName string `yaml:"set-name"`
// 				} `yaml:"eth0"`
// 			} `yaml:"ethernets"`
// 		} `yaml:"network"`
// 	}{}

// 	if err := yaml.Unmarshal(data, &conf); err != nil {
// 		logrus.Error("Ubuntu err ", err)
// 		return err
// 	}

// 	resConf := struct {
// 		NetWork struct {
// 			Ethernets struct {
// 				Ens160 struct {
// 					Addresses   string `yaml:"addresses"`
// 					Dhcp4       string `yaml:"dhcp4"`
// 					Optional    string `yaml:"optional"`
// 					Gateway4    string `yaml:"gateway4"`
// 					Nameservers struct {
// 						Addresses string `yaml:"addresses"`
// 					}
// 				} `yaml:"ens160"`
// 			} `yaml:"ethernets"`
// 			Version  int    `yaml:"version"`
// 			Renderer string `yaml:"renderer"`
// 		} `yaml:"network"`
// 	}{}
// 	resConf.NetWork.Ethernets.Ens160
// 	return nil
// }
