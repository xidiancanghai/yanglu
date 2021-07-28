package model

import (
	"errors"
	"time"
	"yanglu/service/data"

	"github.com/sirupsen/logrus"
)

type Order struct {
	Id           int
	Uid          int
	NodeNum      int
	MonthNum     int
	DiscountCode string
	Status       int
	Total        int // 价格，单位分
	OutTradeNo   string
	UpdateTime   int64
	CreateTime   int64
}

func NewEmptyOrder() *Order {
	return &Order{}
}

func (o *Order) TableName() string {
	return "order_info"
}

func (o *Order) Create() error {
	if o.Uid == 0 || o.Total == 0 || o.OutTradeNo == "" {
		return errors.New("必要参数为空")
	}
	o.CreateTime = time.Now().Unix()
	o.UpdateTime = o.CreateTime
	tx := data.GetDB().Create(o)
	if tx.Error != nil {
		logrus.Error("Create err ", tx.Error)
	}
	return tx.Error
}
