package service

import (
	"errors"
	"yanglu/service/model"

	"github.com/sirupsen/logrus"
)

type UserService struct {
}

func NewUserService() *UserService {
	return &UserService{}
}

func (us *UserService) AddUser(name string, passwd string, authority int, department string) (*model.User, error) {
	u := model.NewUser()
	u, err := u.GetUserByName(name)
	if err != nil {
		logrus.Error("AddUser err ", err)
		return nil, err
	}
	if u.Uid != 0 {
		return nil, errors.New("该用户已存在")
	}
	u.Name = name
	u.Passwd = passwd
	u.Authority = authority
	u.Department = department
	err = u.Create()
	if err != nil {
		logrus.Error("AddUser err ", err)
		return nil, err
	}
	return u, nil
}

func (us *UserService) Login(name string, passwd string) (*model.User, error) {
	u := model.NewUser()
	u, err := u.GetUserByName(name)
	if err != nil {
		logrus.Error("AddUser err ", err)
		return nil, err
	}
	if u.Uid == 0 {
		return nil,errors.New("该用户不存在")
	}
	if 
}
