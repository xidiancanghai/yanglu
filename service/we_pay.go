package service

import (
	"bytes"
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"runtime"
	"strconv"
	"strings"
	"time"
	"yanglu/helper"

	"github.com/sirupsen/logrus"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/native"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
)

var (
	UrlNative                  = "https://api.mch.weixin.qq.com/v3/pay/transactions/native"
	MchId                      = "1611260465"
	AppId                      = "wxf6e86951473e2d1b"
	MchCertificateSerialNumber = "3CE385BA13BA62DC83655DB1DAC7BE51FD25FADB"
)

type WxPayService struct {
	ctx    context.Context
	client *core.Client
}

func NewWxPayService() *WxPayService {
	ws := &WxPayService{}
	err := ws.InitClient()
	if err != nil {
		return nil
	}
	return ws
}

func (ws *WxPayService) GetPrivateKey() (*rsa.PrivateKey, error) {
	_, filename, _, _ := runtime.Caller(0)
	ss := strings.Split(filename, "yanglu")
	absolutelyPath := ss[0] + "yanglu/wx_pay_config/apiclient_key.pem"
	privateKey, err := utils.LoadPrivateKeyWithPath(absolutelyPath)
	if err != nil {
		logrus.Error("GetPrivateKey err ", err)
	}
	return privateKey, nil
}

func (ws *WxPayService) GetPrivateKeyStr() string {
	_, filename, _, _ := runtime.Caller(0)
	ss := strings.Split(filename, "yanglu")
	absolutelyPath := ss[0] + "yanglu/wx_pay_config/apiclient_key.pem"
	data, err := ioutil.ReadFile(absolutelyPath)
	if err != nil {
		logrus.Error("GetPrivateKeyStr err ", err)
		return ""
	}
	return string(data)
}

func (ws *WxPayService) GetCertificate() (*x509.Certificate, error) {
	_, filename, _, _ := runtime.Caller(0)
	ss := strings.Split(filename, "yanglu")
	absolutelyPath := ss[0] + "yanglu/wx_pay_config/apiclient_cert.pem"
	certificate, err := utils.LoadCertificateWithPath(absolutelyPath)
	if err != nil {
		logrus.Error("GetCertificate err ", err)
		return nil, err
	}
	return certificate, nil
}

func (ws *WxPayService) InitClient() error {
	ws.ctx = context.Background()
	privateKey, err := ws.GetPrivateKey()
	if err != nil {
		logrus.Error("InitClient key err ", err)
		return err
	}
	wechatPayCertificate, err := ws.GetCertificate()
	if err != nil {
		logrus.Error("InitClient cert err ", err)
		return err
	}
	certList := []*x509.Certificate{wechatPayCertificate}
	httpClient := new(http.Client)
	opts := []core.ClientOption{
		option.WithWechatPayAuthCipher(MchId, MchCertificateSerialNumber, privateKey, certList),
		option.WithHTTPClient(httpClient),
	}
	ws.client, err = core.NewClient(ws.ctx, opts...)
	if err != nil {
		logrus.Error("InitClient err ", err)
	}
	return nil
}

func (ws *WxPayService) PrePay1(total int, outTradeNo string) (string, error) {

	m := map[string]interface{}{
		"amount": map[string]interface{}{
			"currency": "CNY",
			"total":    total,
		},
		"appid":        AppId,
		"attach":       "自定义数据说明",
		"description":  "引力云测试",
		"mchid":        MchId,
		"notify_url":   "http://matrix.ylysec.com:8080/order/wx_pay_notify",
		"out_trade_no": outTradeNo,
	}

	js, _ := json.Marshal(m)

	client := &http.Client{}

	req, err := http.NewRequest("POST", UrlNative, bytes.NewBuffer(js))

	if err != nil {
		logrus.Error("pre_pay err = ", err)
		return "", err
	}

	req.Header.Add("Content-Type", "application/json")

	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36")

	req.Header.Add("Accept", "application/json")

	token := ws.GetToken("POST", UrlNative, string(js))

	logrus.Info("token = ", token)

	req.Header.Add("Authorization", token)

	rsp, err := client.Do(req)

	if err != nil {
		logrus.Error("pre_pay err = ", err)
		return "", err
	}

	defer rsp.Body.Close()

	data, err := ioutil.ReadAll(rsp.Body)

	logrus.Info("data = ", string(data), " err = ", err)
	if err != nil {
		logrus.Error("pre_pay err = ", err)
		return "", err
	}

	ss := struct {
		CodeUrl string `json:"code_url"`
	}{}

	err = json.Unmarshal(data, &ss)

	if err != nil {
		return "", err
	}

	return ss.CodeUrl, nil
}

func (ws *WxPayService) Sign(method string, url string, nonceStr string, t int64, body string) (string, error) {
	data := ""
	if method == "GET" {
		data = "GET\n%s\n%d\n%s\n\n"
		data = fmt.Sprintf(data, url, t, nonceStr)
	} else {
		data = "POST\n%s\n%d\n%s\n%s\n"
		data = fmt.Sprintf(data, url, t, nonceStr, body)
	}
	hash := crypto.SHA256
	shaNew := hash.New()
	shaNew.Write([]byte(data))
	hashed := shaNew.Sum(nil)

	privateKey, err := ws.GetPrivateKey()
	if err != nil {
		return "", err
	}

	sign, err := rsa.SignPKCS1v15(rand.Reader, privateKey, hash, hashed)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(sign), nil
}

func (ws *WxPayService) GetToken(method string, urlpath string, body string) string {

	t := time.Now().Unix()

	nonceStr := helper.GetRandomStr(20) + strconv.FormatInt(time.Now().UnixNano(), 10)

	url, _ := url.ParseRequestURI(urlpath)

	sign, _ := ws.Sign(method, url.Path, nonceStr, t, body)

	token := fmt.Sprintf("WECHATPAY2-SHA256-RSA2048  mchid=\"%s\",nonce_str=\"%s\",timestamp=\"%d\",serial_no=\"%s\",signature=\"%s\"",
		MchId, nonceStr, t, MchCertificateSerialNumber, sign)
	return token
}

func (ws *WxPayService) PrePay() {
	svc := native.NativeApiService{Client: ws.client}
	resp, result, err := svc.Prepay(ws.ctx,
		native.PrepayRequest{
			Appid:       core.String(AppId),
			Mchid:       core.String(MchId),
			Description: core.String("引力云测试"),
			OutTradeNo:  core.String("12345678910"),
			//TimeExpire:  core.Time(time.Now()),
			Attach:    core.String("自定义数据说明"),
			NotifyUrl: core.String("http://matrix.ylysec.com:8080/order/wx_pay_notify"),
			Amount: &native.Amount{
				Currency: core.String("CNY"),
				Total:    core.Int32(1),
			},
		})
	fmt.Println("resp ", resp)
	data, _ := ioutil.ReadAll(result.Request.Body)
	fmt.Println("result ", string(data))
	fmt.Println("err ", err)
}
