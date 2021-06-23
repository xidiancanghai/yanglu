package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
	"yanglu/service/data"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const (
	AuthoritySuperAdmin         = 1
	AuthorityAddHost            = 2
	AuthorityCheckSoft          = 3
	AuthorityCreateSecurityTask = 4
	AuthorityAddNewUser         = 5
	AuthorityCheckLog           = 6
	AuthorityDeleteUser         = 7
	AuthorityCreateSmartTask    = 8
	AuthorityUpdateSoft         = 9
	AuthorityCreateUserGroup    = 10
	AuthorityCheckDocker        = 11
)

type Ints []int

func (c Ints) Value() (driver.Value, error) {
	b, err := json.Marshal(c)
	return string(b), err
}

func (c *Ints) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), c)
}

type User struct {
	Uid        int    `gorm:"uid" json:"uid"`
	Name       string `gorm:"name" json:"name"`
	Passwd     string `gorm:"passwd" json:"-"`
	Authority  Ints   `gorm:"authority" json:"authority"`
	Department string `gorm:"department" json:"department"`
	IsDelete   int    `gorm:"is_delete" json:"-"`
	UpdateTime int64  `gorm:"update_time" json:"-"`
	CreateTime int64  `gorm:"create_time" json:"-"`
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
	u.CreateTime = time.Now().Unix()
	u.UpdateTime = u.CreateTime
	tx := data.GetDB().Create(u)
	if tx.Error != nil {
		logrus.Error("Create err ", tx)
	}
	return tx.Error
}

func (u *User) GetUserByName(name string) (*User, error) {
	if name == "" {
		return nil, errors.New("参数错误")
	}
	user := new(User)
	tx := data.GetDB().Where(" name = ? and is_delete = 0", name).First(user)
	if tx.Error != nil && tx.Error != gorm.ErrRecordNotFound {
		logrus.Error("GetUserByName err ", tx)
		return nil, tx.Error
	}
	return user, nil
}

func (u *User) GetUserByNamePassWd(name string, passWd string) (*User, error) {
	if name == "" {
		return nil, errors.New("参数错误")
	}
	user := new(User)
	tx := data.GetDB().Where(" name = ? and passwd = ? and is_delete = 0", name, passWd).First(user)
	if tx.Error != nil && tx.Error != gorm.ErrRecordNotFound {
		logrus.Error("GetUserByName err ", tx)
		return nil, tx.Error
	}
	return user, nil
}

func (u *User) GetUserById(uid int) (*User, error) {
	if uid <= 0 {
		return nil, errors.New("参数错误")
	}
	user := new(User)
	tx := data.GetDB().Where(" uid = ? and is_delete = 0", uid).First(user)
	if tx.Error != nil && tx.Error != gorm.ErrRecordNotFound {
		logrus.Error("GetUserById err ", tx)
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
