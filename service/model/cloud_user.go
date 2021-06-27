package model

import (
	"errors"

	"github.com/sirupsen/logrus"
)

type CloudUser struct {
	Uid        int    `db:"uid"`
	Company    string `db:"company"`
	Phone      string `db:"phone"`
	Email      string `db:"email"`
	PassWd     string `db:"passwd"`
	Authority  Ints   `db:"authority"`
	CreateTime int64  `db:"create_time"`
}

func NewCloudUser() *CloudUser {
	return &CloudUser{}
}

func (cu *CloudUser) TableName() string {
	return "cloud_user_info"
}

func (cu *CloudUser) Create() error {
	if cu.Company == "" || cu.Phone == "" || cu.Email == "" || cu.PassWd == "" {
		return errors.New("参数错误")
	}
	err := cu.Create()
	if err != nil {
		logrus.Error("Create err ", err)
		return err
	}
	return nil
}
