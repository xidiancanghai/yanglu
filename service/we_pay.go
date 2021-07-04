package service

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"strings"

	"github.com/sirupsen/logrus"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/h5"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
)

var (
	MchId                      = "1611260465"
	AppId                      = "wx741873479f1ea79b"
	MchCertificateSerialNumber = "3CE385BA13BA62DC83655DB1DAC7BE51FD25FADB"
)

type WxPayService struct {
	ctx    context.Context
	client *core.Client
}

func NewWxPayService() (*WxPayService, error) {
	ws := &WxPayService{}
	err := ws.InitClient()
	if err != nil {
		return nil, err
	}
	return ws, nil
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

func (ws *WxPayService) PrePay() {
	svc := h5.H5ApiService{Client: ws.client}
	resp, result, err := svc.Prepay(ws.ctx, h5.PrepayRequest{
		Appid:         core.String(AppId),
		Mchid:         core.String(MchId),
		Description:   core.String("引力云"),
		OutTradeNo:    core.String("21775251233401201407033233368018"),
		TimeExpire:    core.Time(time.Now()),
		Attach:        core.String("自定义数据说明"),
		NotifyUrl:     core.String("http://www.baidu.com"),
		GoodsTag:      core.String("WXG"),
		LimitPay:      []string{"LimitPay_example"},
		SupportFapiao: core.Bool(false),
		Amount: &h5.Amount{
			Currency: core.String("CNY"),
			Total:    core.Int32(100),
		},
		Detail: &h5.Detail{
			CostPrice: core.Int32(1),
			GoodsDetail: []h5.GoodsDetail{
				{
					GoodsName:        core.String("iphone"),
					MerchantGoodsId:  core.String("ABC"),
					Quantity:         core.Int32(1),
					UnitPrice:        core.Int32(8),
					WechatpayGoodsId: core.String("1001"),
				},
			},
			InvoiceId: core.String("wx123"),
		},
		SceneInfo: &h5.SceneInfo{
			DeviceId: core.String("013467007045764"),
			H5Info: &h5.H5Info{
				AppName:     core.String("王者荣耀"),
				AppUrl:      core.String("https://pay.qq.com"),
				BundleId:    core.String("com.tencent.wzryiOS"),
				PackageName: core.String("com.tencent.tmgp.sgame"),
				Type:        core.String("iOS"),
			},
			PayerClientIp: core.String("14.23.150.211"),
			StoreInfo: &h5.StoreInfo{
				Address:  core.String("广东省深圳市南山区科技中一道10000号"),
				AreaCode: core.String("440305"),
				Id:       core.String("0001"),
				Name:     core.String("腾讯大厦分店"),
			},
		},
		SettleInfo: &h5.SettleInfo{
			ProfitSharing: core.Bool(false),
		},
	})
	fmt.Println("res = ", resp)
	fmt.Println("result = ", result.Response)
	fmt.Println("err = ", err)
}

func (ws *WxPayService) CreateOrder() {

}
