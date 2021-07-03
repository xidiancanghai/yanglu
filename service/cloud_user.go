package service

import (
	"errors"
	"fmt"
	"net/smtp"
	"strings"
	"yanglu/config"
	"yanglu/helper"
	"yanglu/service/model"

	"github.com/jordan-wright/email"
	"github.com/sirupsen/logrus"
)

type CloudUserService struct {
	u *model.CloudUser
}

func NewEmptyCloudUserService() *CloudUserService {
	return &CloudUserService{u: model.NewCloudUser()}
}

func (cs *CloudUserService) Register(company, phone, email, passwd string) (*model.CloudUser, error) {

	if !helper.VerifyMobileFormat(phone) {
		return nil, errors.New("请检查手机号")
	}
	if !helper.VerifyEmailFormat(email) {
		return nil, errors.New("请检查邮箱")
	}

	u, _ := cs.u.GetUser(map[string]interface{}{
		"phone": phone,
	})
	if u.Uid != 0 {
		return nil, errors.New("该手机号已经被注册")
	}
	u, _ = cs.u.GetUser(map[string]interface{}{
		"email": email,
	})
	if u.Uid != 0 {
		return nil, errors.New("该邮箱已经被注册")
	}
	cs.u.Company = company
	cs.u.Phone = phone
	cs.u.Email = email
	cs.u.Passwd = passwd
	cs.u.Authority = model.Ints{}

	tmpPassWd := helper.GetRandomStr(10)
	msg := fmt.Sprintf("引力云的临时密码:%s，请注意查收", tmpPassWd)
	if err := cs.SendEmail(msg); err != nil {
		logrus.Error("Register err ", err)
		return nil, err
	}
	if err := cs.u.Create(); err != nil {
		logrus.Error("Register err ", err)
		return nil, err
	}
	tmp := &model.UserTempPasswd{
		Uid:      cs.u.Uid,
		PassWd:   tmpPassWd,
		IsDelete: 0,
	}
	if err := tmp.Create(); err != nil {
		logrus.Error("Register err ", err)
		return nil, err
	}
	return cs.u, nil
}

func (cs *CloudUserService) SendEmail(msg string) error {
	emailConf := config.GetEmailConf()
	if emailConf == nil {
		return errors.New("邮件配置错误")
	}
	e := email.NewEmail()
	//设置发送方的邮箱
	e.From = emailConf.User
	// 设置接收方的邮箱
	e.To = []string{cs.u.Email}
	//设置主题
	e.Subject = "引力云临时账号密码"
	//设置文件发送的内容
	e.Text = []byte(msg)
	err := e.Send(strings.TrimSpace(emailConf.Addr), smtp.PlainAuth("", emailConf.User, emailConf.Passwd, strings.TrimSpace(emailConf.Host)))
	if err != nil {
		logrus.Error("err = ", err)
		return err
	}
	return nil
}

func (cs *CloudUserService) Login(name string, passwd string) (*model.CloudUser, error) {
	m := map[string]interface{}{}
	if helper.VerifyEmailFormat(name) {
		m["email"] = name
	} else if helper.VerifyMobileFormat(name) {
		m["phone"] = name
	} else {
		return nil, errors.New("请检查手机号或邮箱地址")
	}
	var err error = nil
	cs.u, err = model.NewCloudUser().GetUser(m)
	if err != nil {
		logrus.Error("Login err ", err)
		return nil, err
	}
	if cs.u.Uid == 0 {
		return nil, errors.New("该用户不存在")
	}
	tmpPassWd, _ := model.NewUserTempPasswd().GetPassWd(cs.u.Uid)
	if cs.u.Passwd != passwd && tmpPassWd != passwd {
		return nil, errors.New("账号密码错误")
	}

	return cs.u, nil
}

func (cs *CloudUserService) ResetPasswd(uid int, passwd string) error {
	u, err := model.NewCloudUser().GetUser(map[string]interface{}{
		"uid": uid,
	})
	if err != nil {
		logrus.Error("ResetPasswd err = ", err)
		return err
	}
	if u.Uid == 0 {
		return errors.New("该用户不存在")
	}
	err = u.Updates(map[string]interface{}{
		"passwd": passwd,
	})
	if err != nil {
		logrus.Error("ResetPasswd err = ", err)
		return err
	}
	return nil
}

func (cs *CloudUserService) FindPassWd(account string) (*model.CloudUser, error) {
	var err error = nil
	if helper.VerifyMobileFormat(account) {
		cs.u, err = cs.u.GetUser(map[string]interface{}{
			"phone": account,
		})
	} else if helper.VerifyEmailFormat(account) {
		cs.u, err = cs.u.GetUser(map[string]interface{}{
			"email": account,
		})
	} else {
		return nil, errors.New("请输入正确的邮箱或者手机号")
	}
	if err != nil {
		logrus.Error("FindPassWd err ", err)
		return nil, err
	}
	tmpPassWd := helper.GetRandomStr(10)
	msg := fmt.Sprintf("引力云的临时密码:%s，请注意查收", tmpPassWd)
	if err := cs.SendEmail(msg); err != nil {
		logrus.Error("FindPassWd err ", err)
		return nil, err
	}
	tu, err := model.NewUserTempPasswd().GetTempUser(cs.u.Uid)
	if err != nil {
		logrus.Error("FindPassWd err ", err)
		return nil, err
	}
	tu.PassWd = tmpPassWd
	err = tu.Updates(map[string]interface{}{
		"pass_wd": tu.PassWd,
	})
	if err != nil {
		logrus.Error("FindPassWd err ", err)
		return nil, err
	}
	return cs.u, nil
}
