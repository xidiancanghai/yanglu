package model

import (
	"errors"
	"time"
	"yanglu/service/data"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const (
	AuthorityAddUser = 0
)

type User struct {
	Uid        int    `db:"uid"`
	Name       string `db:"name'`
	Passwd     string `db:"passwd'`
	Authority  int    `db:"authority"`
	Department string `db:"department"`
	UpdateTime int64  `db:"update_time"`
	CreateTime int64  `db:"create_time"`
}

func NewUser() *User {
	return &User{}
}

func (u *User) TableName() string {
	return "user_info"
}

func (u *User) Create() error {
	if u.Name == "" || u.Passwd == "" {
		return errors.New("参数错误")
	}
	if u.CreateTime == 0 {
		u.CreateTime = time.Now().Unix()
		u.UpdateTime = u.CreateTime
	}
	db := data.GetDB()
	tx := db.Create(u)
	if tx.Error != nil {
		logrus.Error("insert error", tx)
	}
	return tx.Error
}

func (u *User) GetUserByName(name string) (*User, error) {
	if name == "" {
		return nil, errors.New("参数错误")
	}
	user := new(User)
	tx := data.GetDB().Where(" name = ?", name).First(user)
	if tx.Error != nil && tx.Error != gorm.ErrRecordNotFound {
		logrus.Error("GetUserByName err ", tx)
		return nil, tx.Error
	}
	return user, nil
}

func (u *User) Updates(m map[string]interface{}) error {
	if u.Uid == 0 {
		return errors.New("参数错误")
	}
	tx := data.GetDB().Model(u).Where(" uid = ?", u.Uid).Updates(m)
	if tx.Error != nil {
		logrus.Error("Updates err", tx)
		return tx.Error
	}
	return nil
}
