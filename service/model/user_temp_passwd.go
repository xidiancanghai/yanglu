package model

import (
	"errors"
	"time"
	"yanglu/service/data"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserTempPasswd struct {
	Uid        int
	PassWd     string
	IsDelete   int
	UpdateTime int64
	CreateTime int64
}

func NewUserTempPasswd() *UserTempPasswd {
	return &UserTempPasswd{}
}

func (up *UserTempPasswd) TableName() string {
	return "user_tmp_passwd"
}

func (up *UserTempPasswd) Create() error {
	if up.Uid == 0 || up.PassWd == "" {
		return errors.New("参数错误")
	}
	up.CreateTime = time.Now().Unix()
	up.UpdateTime = up.CreateTime
	tx := data.GetDB().Create(up)
	if tx.Error != nil {
		logrus.Error("Create err ", tx)
		return tx.Error
	}
	return nil
}

func (up *UserTempPasswd) GetPassWd(uid int) (string, error) {
	if uid == 0 {
		return "", errors.New("参数错误")
	}
	u := new(UserTempPasswd)
	tx := data.GetDB().Where(" uid = ? and is_delete = 0", uid).First(u)
	if tx.Error != nil && tx.Error != gorm.ErrRecordNotFound {
		logrus.Error("GetPassWd err ", tx)
		return "", tx.Error
	}
	return u.PassWd, nil
}

func (up *UserTempPasswd) GetTempUser(uid int) (*UserTempPasswd, error) {
	if uid == 0 {
		return nil, errors.New("参数错误")
	}
	u := new(UserTempPasswd)
	tx := data.GetDB().Where(" uid = ? and is_delete = 0", uid).First(u)
	if tx.Error != nil && tx.Error != gorm.ErrRecordNotFound {
		logrus.Error("GetPassWd err ", tx)
		return nil, tx.Error
	}
	return u, nil
}

func (up *UserTempPasswd) Updates(m map[string]interface{}) error {
	if len(m) == 0 {
		return errors.New("参数错误")
	}
	m["update_time"] = time.Now().Unix()
	tx := data.GetDB().Model(up).Where(" uid = ?", up.Uid).Updates(m)
	if tx.Error != nil {
		logrus.Error("Updates err", tx)
		return tx.Error
	}
	return nil
}
