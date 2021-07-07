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

func (cu *CloudUser) Updates(m map[string]interface{}) error {
	if len(m) == 0 {
		return errors.New("参数错误")
	}
	if cu.Uid == 0 {
		return errors.New("主键id错误")
	}
	tx := data.GetDB().Model(cu).Updates(m)
	if tx.Error != nil {
		logrus.Error("Updates err ", tx)
		return tx.Error
	}
	return nil
}

func (u *CloudUser) GetUserNames(uids []int) map[int]string {
	if len(uids) == 0 {
		return map[int]string{}
	}
	rows, err := data.GetDB().Model(u).Raw("select uid, phone from cloud_user_info where uid in (?)", uids).Rows()
	if err != nil && err != gorm.ErrRecordNotFound {
		logrus.Error("GetUserNames err = ", err)
		return map[int]string{}
	}
	if err != nil {
		return map[int]string{}
	}
	defer rows.Close()
	res := map[int]string{}
	for rows.Next() {
		var uid int
		var name string
		rows.Scan(&uid, &name)
		res[uid] = name
	}
	return res
}
