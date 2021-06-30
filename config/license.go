package config

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"time"
)

type LicenseInfo struct {
	NodeMax       int    `json:"node_max"`
	SmartTask     int    `json:"smart_task"`
	LogManage     int    `json:"log_manage"`
	UserManage    int    `json:"user_manage"`
	UpdateSoft    int    `json:"update_soft"`
	ExpireTimeStr string `json:"expire_time"`
	Code          string `json:"code"`
	IsCloud       int    `json:"is_cloud"`
	ExpireTime    int64
}

var LicenseInfoConf *LicenseInfo

func GetLicensePath() string {
	_, filename, _, _ := runtime.Caller(0)
	return path.Dir(filename)
}

func RSA_decrypter(cipherText []byte) []byte {
	//同加密时，先将私钥从文件中取出，进行二次解码

	privateKey := `-----BEGIN privateKey-----
MIIEpgIBAAKCAQEA+cuo69XmMTnuf9gCb05jR130/O0D2GuYp2EvCmfGxYgyBF0M
VZF+l7hi4ZXNC5j3pbpSAsSqY+rdyVls6XOmibANcXwM2gJFRT/eiU3cAfCVGG+Q
Hw26VP5CdrWyvzN4AVFtwat30H2XZM9fZjKzQU+jdOO4F51EymWBWAO12yQKL09x
iVevgIVKZCoIHt08Fo60D1BhYkF9XlNvUB4O8N0VZ7EIfRp7iIvb+NyuqPksPYTN
h8vp/fYomU2aOXeFO+ur6ENIJSffBJ/5Jg1tiT8pl1/esI9vt1gp6hVk6ehV0HVU
0Hed7dihuIN+iQ+rzFp3EAx3+po6PpWI8aTRXQIDAQABAoIBAQCDCExB+eBYbjUn
pt4QGSSG6Uh77LCZtEdNdfNBGBLyx3DFTTINKNcpKzXjc7stnD48FzKx9QwUgb15
w4kYJOx1SX6jyxjKu2ou20QW7lq8QcFUIiGNyJu3vfCLaOCs34jR2C8k6pddLKRf
40/9wha56daqmFGXQY42RWzGCJw6qPHH6BOxv5nN7ZrcFFJkv1ut5XoDV6coQ1Qm
uojtCDibFLA0NPQFSVs3dWQoprL5iP0PuxiJgB2SlYhafJz0JkZbLYSVE076LKs1
TfX9A3HkZqfrIwiCsOwgyTe8ccEQfav2aTLZ18Z5tQIZrEmzr57s0dJJ8smpY42Q
+Sh/HrHBAoGBAP5cDtouaa8O87bhKLTEOPNKiblYIg8kS+oNAjWqfFOA9HEqnvie
3ppU0OWfNDP+KRK0RZI1s4QFUaVIXjBynFDigs1Oj3eiwnyQW8959sld6A5omQLD
aWey9KY4jSMQf96Wrau7+0tlq9aEPHZOUf/EDzQ6HZ44SKLnsw5l+lWZAoGBAPto
ERHiGyWAZE9NkwXP1qo/vwqV7UOLA2h/x3NJfl0It/2bM7+yK1/NsTH3sug//7Gx
7yYUFwFY54K8ND0yFBfJwBYO9S44tPa7RL+24+h8Wv2uRVH8t1p/hByifm5sQ89g
0kUb2maRNEG4Nr5b8WgWDL7xW7zwdr4h451RuOxlAoGBAN9uPh2pZ2w+1hZu8jLX
qq8Laq9Em6ZxW73cb3R0rYEsZYyCw+Hgq+klV9pcXp7NQ4cKk5I7WdU+Vtm5GAO2
qVjjn8J4n82XSv/gTsWRfRalRpJhjwa+YHiW2hCo3FFQzQHZmfSgy70PMkLdJb9G
TQE+V5yo8BtrlWKD5OVg7uUBAoGBALxr7I84pRvgNE3zWAN9F+gVSrtGTsLWh0Fk
szlGlbxocm1SSuCbSkLNSMaoXKAQhzTIpPInEG7kEO882vIafFCv1pUxtLRlr8tT
41zEV9/Ag/mpCS8/drQO6hcD+joMTm44pJ9DFtbUcRPoWq8k/bMDA56ACj97AaPS
hmM8Apq5AoGBAMuTf5hwqhCoJmfKQ1uyydRySGl3O/lF3Vam2D/h3kETI+ebZ0nI
1QLz0bObt5YgXZl7nmnW4dlxBMWRvpyF54afG7X7F2ddDgwF7tydGMLb/+XpYEXJ
F1VY5dv1sasQfsINXbxPojAxAe/hlZgZOhWpKtB06Db7KStCpBv8pDGe
-----END privateKey-----`
	block, _ := pem.Decode([]byte(privateKey))
	PrivateKey, _ := x509.ParsePKCS1PrivateKey(block.Bytes)
	//二次解码完毕，调用解密函数
	afterDecrypter, _ := rsa.DecryptPKCS1v15(rand.Reader, PrivateKey, cipherText)
	return afterDecrypter
}

func InitLicenseConfig() error {
	filePath := GetLicensePath() + "/license.txt"
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		LicenseInfoConf = new(LicenseInfo)
		LicenseInfoConf.NodeMax = 5
		LicenseInfoConf.SmartTask = 0
		LicenseInfoConf.LogManage = 0
		LicenseInfoConf.UserManage = 0
		LicenseInfoConf.UpdateSoft = 0
		LicenseInfoConf.ExpireTime = time.Now().Unix() + 3600*24*365*100
		return nil
	}

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	bytes, _ = hex.DecodeString(string(bytes))

	bytes = RSA_decrypter(bytes)

	LicenseInfoConf = new(LicenseInfo)

	err = json.Unmarshal(bytes, LicenseInfoConf)

	if err != nil {
		return err
	}

	loc, _ := time.LoadLocation("Local")

	theTime, err := time.ParseInLocation("2006-01-02", LicenseInfoConf.ExpireTimeStr, loc)

	if err != nil {
		return err
	}

	LicenseInfoConf.ExpireTime = theTime.Unix()
	return nil
}

func IsCloud() bool {
	if LicenseInfoConf == nil {
		return false
	}
	return LicenseInfoConf.IsCloud == 1
}
