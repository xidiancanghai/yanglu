package model

import (
	"errors"
	"fmt"
	"time"
	"yanglu/service/data"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ActionLog struct {
	Id         int    `json:"id"`
	Uid        int    `json:"uid"`
	Detail     string `json:"detail"`
	CreateTime int64  `json:"create_time"`
}

func NewActionLog() *ActionLog {
	return &ActionLog{}
}

func (al *ActionLog) TableName() string {
	return "action_log"
}

func (a *ActionLog) Create() error {
	if a.Uid == 0 {
		return errors.New("参数错误")
	}
	a.CreateTime = time.Now().Unix()
	db := data.GetDB()
	tx := db.Create(a)
	if tx.Error != nil {
		logrus.Error("insert error", tx)
	}
	return tx.Error
}

func (a *ActionLog) List(lastId int) ([]*ActionLog, error) {
	sqll := ""
	if lastId == -1 {
		sqll = "select * from action_log order by id desc limit 50"
	} else {
		sqll = fmt.Sprintf("select * from action_log where id < %d order by id desc limit 50", lastId)
	}
	list := []*ActionLog{}
	tx := data.GetDB().Model(a).Raw(sqll).Find(&list)
	if tx.Error != nil && tx.Error != gorm.ErrRecordNotFound {
		logrus.Error("List err ", tx)
		return nil, tx.Error
	}
	return list, nil
}
