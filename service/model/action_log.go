package model

import (
	"errors"
	"fmt"
	"time"
	"yanglu/config"
	"yanglu/service/data"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const (
	ActionTypeLogin   = 1
	ActionTypeAddHost = 2
	ActionTypeAddTask = 3
	ActionTypeAddUser = 4
	ActionTypeLogout  = 5
)

type ActionLog struct {
	Id         int    `json:"id"`
	Uid        int    `json:"uid"`
	Type       int    `json:"type"`
	Ip         string `json:"ip"`
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
	if config.LicenseInfoConf.LogManage == 0 {
		return errors.New("当前系统无日志管理权限")
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

func (a *ActionLog) SearchLog(startTime int, endTime int, action int, ip string) ([]*ActionLog, error) {

	conds := []string{}
	where := ""
	if startTime != 0 {
		conds = append(conds, fmt.Sprintf(" create_time > %d ", startTime))
	}
	if endTime != 0 {
		conds = append(conds, fmt.Sprintf(" create_time < %d ", endTime))
	}
	if action != 0 {
		conds = append(conds, fmt.Sprintf(" type = %d ", action))
	}
	if ip != "" {
		conds = append(conds, `ip like '%`+ip+`%'`)
	}

	if len(conds) != 0 {
		where = " where "
		for k, v := range conds {
			if k < len(conds)-1 {
				where = where + v + " and "
			} else {
				where = where + v
			}
		}
	}

	sqll := "select * from " + a.TableName() + where + " order by id desc limit 100 "

	if config.IsTest() {
		logrus.Info("sqll = ", sqll)
	}

	list := []*ActionLog{}
	tx := data.GetDB().Model(a).Raw(sqll).Find(&list)
	if tx.Error != nil && tx.Error != gorm.ErrRecordNotFound {
		logrus.Error("List err ", tx)
		return nil, tx.Error
	}
	return list, nil
}
