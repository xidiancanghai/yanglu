package service

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"yanglu/config"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
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
	err = nil
	if err != nil {
		logrus.Error("SetIp err ", err)
		return err
	}
	if config.IsLocal() {
		data = []byte(`Ubuntu 20.04.2 LTS \n \l`)
	}
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
		ss.Ubuntu()
	}
	return nil
}

func (ss *SetIpService) Ubuntu() error {
	files := []string{}
	path := "/etc/netplan"
	infos, err := ioutil.ReadDir(path)
	if err != nil {
		logrus.Error("Ubuntu err ", err)
		return err
	}

	for _, x := range infos {
		realPath := path + "/" + x.Name()
		if x.IsDir() {
			continue
		} else {
			files = append(files, realPath)
		}
	}

	if len(files) < 2 {
		return errors.New("当前只有一个网卡，不满足条件")
	}

	filePath := files[0]

	file, err := os.Open(filePath)
	if err != nil {
		logrus.Error("Ubuntu err ", err)
		return err
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		logrus.Error("Ubuntu err ", err)
		return err
	}
	conf := struct {
		Network struct {
			Version   string `yaml:"version"`
			Ethernets struct {
				Eth0 struct {
					Dhcp4 bool `yaml:"dhcp4"`
					Match struct {
						Macaddress string `yaml:"macaddress"`
					} `yaml:"match"`
					SetName string `yaml:"set-name"`
				} `yaml:"eth0"`
			} `yaml:"ethernets"`
		} `yaml:"network"`
	}{}

	if err := yaml.Unmarshal(data, &conf); err != nil {
		logrus.Error("Ubuntu err ", err)
		return err
	}

	logrus.Info("debug = ", conf)
	// resConf := struct {
	// 	NetWork struct {
	// 		Ethernets struct {
	// 			Ens160 struct {
	// 				Addresses   string `yaml:"addresses"`
	// 				Dhcp4       string `yaml:"dhcp4"`
	// 				Optional    string `yaml:"optional"`
	// 				Gateway4    string `yaml:"gateway4"`
	// 				Nameservers struct {
	// 					Addresses string `yaml:"addresses"`
	// 				}
	// 			} `yaml:"ens160"`
	// 		} `yaml:"ethernets"`
	// 		Version  int    `yaml:"version"`
	// 		Renderer string `yaml:"renderer"`
	// 	} `yaml:"network"`
	// }{}

	return nil
}
