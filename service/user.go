package service

import (
	"errors"
	"yanglu/config"
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
	u.Authority = model.Ints{authority}
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
		adminUser, adminPassWd := config.GetAdminInfo()
		if adminUser == name && adminPassWd == passwd {
			u.Authority = model.Ints{model.AuthoritySuperAdmin}
			u.Name = adminUser
			u.Passwd = adminPassWd
			err := u.Create()
			if err != nil {
				return nil, err
			}
			u, _ = model.NewUser().GetUserByNamePassWd(u.Name, u.Passwd)
			return u, nil
		}
	}
	if u.Uid == 0 {
		return nil, errors.New("该用户不存在")
	}
	if u.Passwd != passwd {
		return nil, errors.New("密码错误")
	}
	return u, nil
}

func (us *UserService) UserInfo(uid int) (*model.User, error) {
	u := model.NewUser()
	u, err := u.GetUserById(uid)
	if err != nil {
		logrus.Error("UserInfo err ", err)
		return nil, err
	}
	if u.Uid == 0 {
		return nil, errors.New("该用户不存在")
	}
	return u, nil
}

func (us *UserService) SetAuthority(name string, authority int) error {
	user, err := model.NewUser().GetUserByName(name)
	if err != nil {
		return err
	}
	if user.Uid == 0 {
		return errors.New("该用户不存在")
	}
	user.Authority = append(user.Authority, authority)
	err = user.Updates(map[string]interface{}{
		"authority": user.Authority,
	})
	if err != nil {
		logrus.Error("SetAuthority err ", err)
		return err
	}
	return nil
}

func (us *UserService) GetAdminUser() (*model.User, error) {
	return model.NewUser().GetUserByNamePassWd(config.GetAdminInfo())
}

func (us *UserService) HasAuthorityByUser(user *model.User, authority int) bool {
	adminUser, err := us.GetAdminUser()
	if err != nil {
		return false
	}
	if user.Uid == adminUser.Uid {
		return true
	}
	for _, v := range user.Authority {
		if v == authority {
			return true
		}
	}
	return false
}

func (us *UserService) HasAuthority(uid int, authority int) bool {
	user, err := model.NewUser().GetUserById(uid)
	if err != nil {
		return false
	}
	return us.HasAuthorityByUser(user, authority)
}

func (us *UserService) DeleteUser(name string) error {
	user, err := model.NewUser().GetUserByName(name)
	if err != nil {
		return err
	}
	err = user.Updates(map[string]interface{}{
		"is_delete": 1,
	})
	if err != nil {
		logrus.Error("DeleteUser err ", err)
		return err
	}
	return nil
}
