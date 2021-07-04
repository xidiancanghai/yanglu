package service

import (
	"crypto/rsa"
)

var (
	MahId                      = "1611260465"
	AppId                      = "wx741873479f1ea79b"
	MchCertificateSerialNumber = "3CE385BA13BA62DC83655DB1DAC7BE51FD25FADB"
)

type WxPayService struct {
	privateKey *rsa.PrivateKey
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
