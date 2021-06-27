package controller

import (
	"errors"
	"yanglu/def"
	"yanglu/helper"
	"yanglu/service"
	"yanglu/service/model"

	"github.com/gin-gonic/gin"
)

type User struct {
}

func NewUser() *User {
	return &User{}
}

func (u *User) Register(ctx *gin.Context) {
	params := &struct {
		Company      string `form:"company" binding:"required"`
		Phone        string `form:"phone" binding:"required"`
		email        string `form:"email" binding:"required"`
		passwd       string `form:"passwd" binding:"required"`
		CaptchaId    string `form:"captcha_id" binding:"required"`
		CaptchaValue string `form:"captcha_value" binding:"required"`
	}{}

	if err := ctx.ShouldBind(params); err != nil {
		helper.ErrRsp(ctx, def.CodeErr, err.Error(), err)
		return
	}

	if err := service.NewUtilService().ValidateCaptcha(params.CaptchaId, params.CaptchaValue); err != nil {
		helper.ErrRsp(ctx, def.CodeErr, err.Error(), err)
		return
	}

}

func (u *User) AddUser(ctx *gin.Context) {

	params := &struct {
		Name       string `form:"name" binding:"required"`
		Passwd     string `form:"passwd" binding:"required"`
		Authority  int    `form:"authority" binding:"required"`
		Department string `form:"department" binding:"required"`
	}{}

	if err := ctx.ShouldBind(params); err != nil {
		helper.ErrRsp(ctx, def.CodeErr, err.Error(), err)
		return
	}

	uid := ctx.GetInt("uid")
	us := service.NewUserService()
	user, err := us.UserInfo(uid)
	if err != nil {
		helper.ErrRsp(ctx, def.CodeErr, err.Error(), err)
		return
	}
	if !us.HasAuthorityByUser(user, model.AuthoritySuperAdmin) {
		helper.ErrRsp(ctx, def.CodeErr, "你没有权限添加用户", errors.New("你没有权限添加用户"))
		return
	}

	usr, err := service.NewUserService().AddUser(params.Name, params.Passwd, params.Authority, params.Department)
	if err != nil {
		helper.ErrRsp(ctx, def.CodeErr, err.Error(), err)
		return
	}
	service.NewActionLogService(uid).AddUser(usr)
	helper.OKRsp(ctx, gin.H{})
}

func (u *User) Login(ctx *gin.Context) {
	params := &struct {
		Name         string `form:"name" binding:"required"`
		Passwd       string `form:"passwd" binding:"required"`
		CaptchaId    string `form:"captcha_id"`
		CaptchaValue string `form:"captcha_value"`
	}{}

	if err := ctx.ShouldBind(params); err != nil {
		helper.ErrRsp(ctx, def.CodeErr, err.Error(), err)
		return
	}

	if err := service.NewUtilService().ValidateCaptcha(params.CaptchaId, params.CaptchaValue); err != nil {
		helper.ErrRsp(ctx, def.CodeErr, err.Error(), err)
		return
	}

	usr, err := service.NewUserService().Login(params.Name, params.Passwd)
	if err != nil {
		helper.ErrRsp(ctx, def.CodeErr, err.Error(), err)
		return
	}
	token, err := service.NewTokenService(usr.Uid).BuildToken()
	if err != nil {
		helper.ErrRsp(ctx, def.CodeErr, err.Error(), err)
		return
	}
	service.NewActionLogService(usr.Uid).Login()
	res := gin.H{
		"token": token,
	}
	helper.OKRsp(ctx, res)
}

func (u *User) GetUserInfo(ctx *gin.Context) {
	uid := ctx.GetInt("uid")
	user, err := service.NewUserService().UserInfo(uid)
	if err != nil {
		helper.ErrRsp(ctx, def.CodeErr, err.Error(), err)
		return
	}
	res := gin.H{
		"name":       user.Name,
		"authority":  user.Authority,
		"department": user.Department,
	}
	helper.OKRsp(ctx, res)
}

func (u *User) SetAuthority(ctx *gin.Context) {
	uid := ctx.GetInt("uid")
	params := &struct {
		TargetName string `form:"target_name" binding:"required"`
		Authority  int    `form:"authority" binding:"required"`
	}{}

	if err := ctx.ShouldBind(params); err != nil {
		helper.ErrRsp(ctx, def.CodeErr, err.Error(), err)
		return
	}

	us := service.NewUserService()
	user, err := us.UserInfo(uid)
	if err != nil {
		helper.ErrRsp(ctx, def.CodeErr, err.Error(), err)
		return
	}
	if !us.HasAuthorityByUser(user, model.AuthoritySuperAdmin) {
		helper.ErrRsp(ctx, def.CodeErr, "你没有权限设置他人权限", errors.New("你没有权限设置他人权限"))
		return
	}

	err = service.NewUserService().SetAuthority(params.TargetName, params.Authority)
	if err != nil {
		helper.ErrRsp(ctx, def.CodeErr, err.Error(), err)
		return
	}
	helper.OKRsp(ctx, gin.H{})
}

func (u *User) DeleteUser(ctx *gin.Context) {
	uid := ctx.GetInt("uid")
	params := &struct {
		TargetName string `form:"target_name" binding:"required"`
	}{}

	if err := ctx.ShouldBind(params); err != nil {
		helper.ErrRsp(ctx, def.CodeErr, err.Error(), err)
		return
	}

	us := service.NewUserService()
	user, err := us.UserInfo(uid)
	if err != nil {
		helper.ErrRsp(ctx, def.CodeErr, err.Error(), err)
		return
	}
	if !us.HasAuthorityByUser(user, model.AuthoritySuperAdmin) {
		helper.ErrRsp(ctx, def.CodeErr, "你没有权限设置他人权限", errors.New("你没有权限设置他人权限"))
		return
	}

	err = service.NewUserService().DeleteUser(params.TargetName)
	if err != nil {
		helper.ErrRsp(ctx, def.CodeErr, err.Error(), err)
		return
	}
	helper.OKRsp(ctx, gin.H{})
}

func (u *User) ListUsers(ctx *gin.Context) {
	us := service.NewUserService()
	list, err := us.ListUsers()
	if err != nil {
		helper.ErrRsp(ctx, def.CodeErr, err.Error(), err)
		return
	}
	helper.OKRsp(ctx, gin.H{"list": list})
}
