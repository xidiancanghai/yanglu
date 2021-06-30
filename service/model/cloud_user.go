package model

import (
	"errors"
	"fmt"
	"yanglu/service/data"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CloudUser struct {
	Uid        int    `gorm:"primaryKey"`
	Company    string `gorm:"company"`
	Phone      string `gorm:"phone"`
	Email      string `gorm:"email"`
	Passwd     string `gorm:"passwd"`
	Authority  Ints   `gorm:"authority"`
	CreateTime int64  `gorm:"create_time"`
}

func NewCloudUser() *CloudUser {
	return &CloudUser{}
}

func (cu *CloudUser) TableName() string {
	return "cloud_user_info"
}

func (cu *CloudUser) Create() error {
	if cu.Company == "" || cu.Phone == "" || cu.Email == "" || cu.Passwd == "" {
		return errors.New("参数错误")
	}
	tx := data.GetDB().Create(cu)
	if tx.Error != nil {
		logrus.Error("Create err ", tx)
		return tx.Error
	}
	return nil
}

func (cu *CloudUser) GetUser(cond map[string]interface{}) (*CloudUser, error) {
	if len(cond) == 0 {
		return nil, errors.New("参数错误")
	}
	u := new(CloudUser)
	key := []string{}
	values := []interface{}{}
	conds := ""
	for k, v := range cond {
		key = append(key, k)
		values = append(values, v)
		conds = conds + fmt.Sprintf(" %s = ? ", k)
		if len(key) < len(cond)-1 {
			conds = conds + " and "
		}
	}
	tx := data.GetDB().Model(cu).Where(conds, values...).First(u)
	if tx.Error != nil && tx.Error != gorm.ErrRecordNotFound {
		logrus.Error("Create err ", tx)
		return nil, tx.Error
	}
	return u, nil
}
