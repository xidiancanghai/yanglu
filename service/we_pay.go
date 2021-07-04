package service

import (
	"bytes"
	"crypto/md5"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"yanglu/helper"

	"github.com/sirupsen/logrus"
)

var (
	UrlUnifiedorder            = "https://api.mch.weixin.qq.com/pay/unifiedorder"
	MchId                      = "1611260465"
	AppId                      = "wx741873479f1ea79b"
	MchCertificateSerialNumber = "3CE385BA13BA62DC83655DB1DAC7BE51FD25FADB"
)

type WxPayService struct {
}

// func (ws *WxPayService) InitPrivateKey() error {
// 	_, filename, _, _ := runtime.Caller(0)
// 	ss := strings.Split(filename, "yanglu")
// 	absolutelyPath := ss[0] + "yanglu/wx_pay_config/apiclient_key.pem"
// 	keyBuffer, err := ioutil.ReadFile(absolutelyPath)
// 	if err != nil {
// 		logrus.Error("InitPrivateKey err ", err)
// 		return err
// 	}
// 	block, _ := pem.Decode(keyBuffer)
// 	if block == nil {
// 		return errors.New("private key error")
// 	}
// 	ws.privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (ws *WxPayService) InitClient() {
// 	ctx := context.Background()
// 	httpClient := new(http.Client)
// 	opts := []core.ClientOption{
// 		option.WithWechatPayAuthCipher(MahId, MchCertificateSerialNumber, ws.privateKey),
// 	}
// }

// func NewWxPayService() {
// 	ws := &WxPayService{}
// 	ws.InitPrivateKey()
// 	ws.InitClient()
// }

// func (ws *WxPayService) PrePay() {

// }

func (ws *WxPayService) PrePay() error {
	// request := struct {
	// 	AppId          string `xml:"appid"`
	// 	Attach         string `xml:"attach"`
	// 	Body           string `xml:"body"`
	// 	MchId          string `xml:"mch_id"`
	// 	NonceStr       string `xml:"nonce_str"`
	// 	NotifyUrl      string `xml:"notify_url"`
	// 	Openid         string `xml:"openid"`
	// 	OuTradeNo      string `xml:"out_trade_no"`
	// 	SpbillCreateIp string `xml:"spbill_create_ip"`
	// 	TotalFee       int    `xml:"total_fee"`
	// 	TradeType      string `xml:"trade_type"`
	// 	SceneInfo      struct {
	// 		H5Info struct {
	// 			Type        string `xml:"type"`
	// 			AppName     string `xml:"app_name"`
	// 			PackageName string `xml:"package_name"`
	// 		} `xml:"h5_info"`
	// 	} `xml:"scene_info"`
	// }{}
	// request.AppId = AppId
	// request.Attach = "测试"
	// request.Body = "h5支付测试"
	// request.MchId = MchId
	// request.NonceStr = helper.GetRandomStr(32)
	// request.NotifyUrl = ""
	m := map[string]interface{}{
		"appid":            AppId,
		"mch_id":           MchId,
		"nonce_str":        helper.GetRandomStr(20),
		"body":             "h5测试",
		"out_trade_no":     helper.GetRandomStr(32),
		"total_fee":        1,
		"spbill_create_ip": "127.0.0.1",
		"notify_url":       "http://matrix.ylysec.com:8080/order/wx_notify",
		"trade_type":       "MWEB",
		"scene_info": map[string]interface{}{
			"h5_info": map[string]interface{}{
				"type":     "Wap",
				"wap_url":  "http://matrix.ylysec.com:8080",
				"wap_name": "引力云",
			},
		},
	}
	sign := ws.Sign(m)
	m["sign"] = sign

	buf, _ := xml.Marshal(MapToXml(m))

	client := &http.Client{}
	req, err := http.NewRequest("POST", UrlUnifiedorder, bytes.NewBuffer(buf))

	if err != nil {
		logrus.Error("pre_pay err = ", err)
		return err
	}

	req.Header.Add("Content-Type", "application/xml")

	rsp, err := client.Do(req)

	if err != nil {
		logrus.Error("pre_pay err = ", err)
		return err
	}

	defer rsp.Body.Close()

	data, err := ioutil.ReadAll(rsp.Body)

	if err != nil {
		logrus.Error("pre_pay err = ", err)
		return err
	}

	fmt.Println("data = ", string(data))

	return nil
}

func (ws *WxPayService) Sign(m map[string]interface{}) string {
	list := []string{}
	for k := range m {
		if k == "sign" {
			continue
		}
		list = append(list, k)
	}
	sort.Strings(list)
	bf := bytes.NewBuffer(make([]byte, 0, 400))
	for _, v := range list {
		bf.WriteString(fmt.Sprintf("%s=%v&", v, m[v]))
	}
	bf.WriteString(fmt.Sprintf("key=%s", "isosIm75f386aef0ee6e32b22dylySec"))
	md5 := md5.New()
	md5.Write(bf.Bytes())
	return fmt.Sprintf("%X", md5.Sum(nil))
}
