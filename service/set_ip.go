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

	mapstructure "github.com/goinggo/mapstructure"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type SetIpService struct {
	ip string
}

func NewSetIpService(ip string) *SetIpService {
	return &SetIpService{ip: ip}
}

func (ss *SetIpService) IpHasUsed() (bool, error) {
	command := fmt.Sprintf("ping -c 1 %s > /dev/null && echo true || echo false", ss.ip)
	output, err := exec.Command("/bin/sh", "-c", command).Output()
	if err != nil {
		logrus.Error("IpHasUsed err ", err)
		return false, err
	}
	logrus.WithFields(logrus.Fields{
		"ip":     ss.ip,
		"output": string(output),
		"cmd":    command,
	}).Info("IpHasUsed")
	if strings.Contains(string(output), "true") {
		return true, nil
	}
	return false, nil
}

func (ss *SetIpService) SetIp() error {
	is, err := ss.IpHasUsed()
	if err != nil {
		logrus.Error("SetIp err ", err)
		return err
	}
	if is {
		return errors.New("该ip地址已经被占用")
	}
	cmd := exec.Command("cat", "/etc/issue")
	data, err := cmd.CombinedOutput()
	if config.IsLocal() {
		err = nil
	}
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

func (ss *SetIpService) Gateway4() (string, error) {
	cmd := exec.Command("bash", "-c", `route -n|awk '{if ($1=="0.0.0.0") print $2}'`)
	output, err := cmd.CombinedOutput()
	if config.IsLocal() {
		return "127.0.0.1", nil
	}
	if err != nil {
		logrus.Error("Gateway4 err ", err)
		return "", err
	}
	return string(output), nil
}

func (ss *SetIpService) GetNameServer() (string, error) {
	cmd := exec.Command("bash", "-c", `cat /etc/resolv.conf|awk '/nameserver [0-9]+\.[0-9]+\.[0-9]+\.[0-9]+/ {print $2}'`)
	output, err := cmd.CombinedOutput()
	if err != nil {
		logrus.Error("SetIp err ", err)
		return "", err
	}
	if len(output) == 0 {
		return "", errors.New("")
	}
	return string(output), nil
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

	card := ""
	version := 0
	effectiveFile := ""
	for i := range files {
		effectiveFile = files[i]
		version, card, err = ss.GetNetWorkConf(files[i])
		if version != 0 && card != "" {
			break
		}
	}
	effectiveFile = "/Users/weipai-liuxiang/go/src/yanglu/service/50-cloud-init.yaml"
	if card == "" {
		logrus.Error("Ubuntu err ", err)
		return errors.New("找不到网卡信息")
	}

	gateway4, err := ss.Gateway4()
	if err != nil {
		logrus.Error("Ubuntu err ", err)
		return err
	}
	nameServer, err := ss.GetNameServer()
	if err != nil {
		logrus.Error("Ubuntu err ", err)
		return err
	}

	address := fmt.Sprintf("%s/24", ss.ip)

	// res := map[string]interface{}{
	// 	"network": map[string]interface{}{
	// 		card: map[string]interface{}{
	// 			"dhcp4":    "no",
	// 			"address":  string(addressBytes),
	// 			"optional": true,
	// 			"gateway4": gateway4,
	// 			"nameservers": map[string]interface{}{
	// 				"addresses": string(nameserverBytes),
	// 			},
	// 		},
	// 	},
	// 	"version": version,
	// }

	res := `
network:
    ethernets:
        ens33:
        	dhcp4: no
            	addresses: [{0}/24]
            	optional: true
            	gateway4: {1}
            	nameservers:
                    addresses: [{2}]
 
    version: {3}
	`
	res = strings.ReplaceAll(res, "{0}", address)
	res = strings.ReplaceAll(res, "{1}", gateway4)
	res = strings.ReplaceAll(res, "{2}", nameServer)
	res = strings.ReplaceAll(res, "{3}", strconv.Itoa(version))
	err = ioutil.WriteFile(effectiveFile, []byte(res), 0777)
	if err != nil {
		logrus.Error("Ubuntu err ", err)
		return err
	}
	return nil
}

func (ss *SetIpService) GetNetWorkConf(filePath string) (int, string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		logrus.Error("Ubuntu err ", err)
		return 0, "", err
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		logrus.Error("Ubuntu err ", err)
		return 0, "", err
	}

	conf := map[string]map[string]interface{}{}

	if err := yaml.Unmarshal(data, &conf); err != nil {
		logrus.Error("Ubuntu err ", err)
		return 0, "", err
	}

	netWork, ok := conf["network"]
	if !ok {
		return 0, "", errors.New("配置文件错误")
	}
	ethernets, ok := netWork["ethernets"]
	if !ok {
		return 0, "", errors.New("配置文件错误")
	}
	versionStr, ok := netWork["version"]
	if !ok {
		return 0, "", errors.New("配置文件错误")
	}
	version, _ := InterfaceToInt(versionStr)

	eth := struct {
		Dhcp4 bool `yaml:"dhcp4" json:"dhcp4"`
		Match struct {
			Macaddress string `yaml:"macaddress" json:"macaddress"`
		} `yaml:"match" json:"match"`
		SetName string `yaml:"set-name" json:"set-name"`
	}{}
	//logrus.Info("ethernets ", ethernets)

	eth0, _ := ethernets.(map[interface{}]interface{})

	//logrus.Info("eth = ", eth0, " ok = ", ok)

	ethConf := convert(eth0)

	for k, v := range ethConf {
		mapstructure.Decode(v, &eth)
		if eth.Match.Macaddress != "" {
			return version, k, nil
		}
	}

	return version, "", errors.New("找不到网卡")
}

func convert(m map[interface{}]interface{}) map[string]interface{} {
	res := map[string]interface{}{}
	for k, v := range m {
		switch v2 := v.(type) {
		case map[interface{}]interface{}:
			res[fmt.Sprint(k)] = convert(v2)
		default:
			res[fmt.Sprint(k)] = v
		}
	}
	return res
}

func InterfaceToInt(in interface{}) (int, error) {
	switch in.(type) {
	case string:
		return strconv.Atoi(in.(string))
	case int:
		return in.(int), nil
	case int32:
		return int(in.(int32)), nil
	case int64:
		return int(in.(int64)), nil
	case float32:
		return int(in.(float32)), nil
	case float64:
		return int(in.(float64)), nil
	default:
		logrus.Error("InterfaceToInt in ", in)
	}
	return 0, errors.New("转换错误")
}
