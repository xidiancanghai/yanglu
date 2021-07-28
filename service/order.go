package service

import (
	"strconv"
	"time"
	"yanglu/helper"
	"yanglu/service/model"

	"github.com/sirupsen/logrus"
	"github.com/skip2/go-qrcode"
)

const (
	NodePrice  = 100 // 100fen
	MonthPrice = 100 //100fen
)

type OrderService struct {
	uid   int
	order *model.Order
}

func NewOrderService(uid int) *OrderService {
	return &OrderService{
		uid:   uid,
		order: model.NewEmptyOrder(),
	}
}

func (os *OrderService) Create(nodeNum int, monthNum int, discountCode string) ([]byte, error) {

	totolFeed := nodeNum * NodePrice * monthNum * MonthPrice

	os.order.Uid = os.uid
	os.order.NodeNum = nodeNum
	os.order.MonthNum = monthNum
	os.order.DiscountCode = discountCode
	os.order.Status = 1
	os.order.Total = totolFeed
	os.order.OutTradeNo = helper.GetRandomStr(10) + strconv.FormatInt(time.Now().UnixNano(), 10)

	err := os.order.Create()

	if err != nil {
		logrus.Error("Create err ", err)
		return nil, err
	}

	ws := NewWxPayService()
	codeUrl, err := ws.PrePay1(os.order.Total, os.order.OutTradeNo)

	if err != nil {
		return nil, err
	}

	q, err := qrcode.New(codeUrl, qrcode.Medium)
	if err != nil {
		logrus.Error("Create")
		return nil, err
	}
	png, err := q.PNG(256)
	if err != nil {
		logrus.Error("Create")
	}
	return png, nil
}
