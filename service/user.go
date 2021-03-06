package service

import (
	"errors"
	"sort"
	"yanglu/config"
	"yanglu/service/model"

	"github.com/sirupsen/logrus"
)

type UserService struct {
}

func NewUserService() *UserService {
	return &UserService{}
}

func (us *UserService) AddUser(name string, passwd string, authority []int, department string) (*model.User, error) {

	if config.LicenseInfoConf.UserManage == 0 {
		return nil, errors.New("当前系统没有添加用户权限")
	}

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

func (us *UserService) SetAuthority(name string, authority []int) error {
	user, err := model.NewUser().GetUserByName(name)
	if err != nil {
		return err
	}
	if user.Uid == 0 {
		return errors.New("该用户不存在")
	}
	user.Authority = append(user.Authority, authority...)
	// 去重复
	sort.Ints(user.Authority)
	ret := []int{}
	for i := 0; i < len(user.Authority); {
		ret = append(ret, user.Authority[i])
		j := i
		for j < len(user.Authority) && user.Authority[i] == user.Authority[j] {
			j++
		}
		i = j
	}
	err = user.Updates(map[string]interface{}{
		"authority": model.Ints(ret),
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
		if v == authority || v == model.AuthoritySuperAdmin {
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

func (us *UserService) ListUsers() ([]*model.User, error) {
	return model.NewUser().ListUsers()
}

func (us *UserService) ResetPasswd(uid int, passwd string) error {
	user, err := model.NewUser().GetUserById(uid)
	if err != nil {
		logrus.Error("ResetPasswd err ", err)
		return err
	}
	err = user.Updates(map[string]interface{}{
		"passwd": passwd,
	})
	if err != nil {
		logrus.Error("ResetPasswd err ", err)
		return err
	}
	return nil
}
