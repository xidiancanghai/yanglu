package controller

import (
	"yanglu/def"
	"yanglu/helper"
	"yanglu/service"

	"github.com/gin-gonic/gin"
)

type User struct {
}

func NewUser() *User {
	return &User{}
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
	usr, err := service.NewUserService().AddUser(params.Name, params.Passwd, params.Authority, params.Department)
	if err != nil {
		helper.ErrRsp(ctx, def.CodeErr, err.Error(), err)
		return
	}
	token, err := service.NewTokenService(usr.Uid).BuildToken()
	if err != nil {
		helper.ErrRsp(ctx, def.CodeErr, err.Error(), err)
	}
	res := gin.H{
		"token": token,
	}
	helper.OKRsp(ctx, res)
}

func (u *User) Login(ctx *gin.Context) {
	params := &struct {
		Name   string `form:"name" binding:"required"`
		Passwd string `form:"passwd" binding:"required"`
	}{}

	if err := ctx.ShouldBind(params); err != nil {
		helper.ErrRsp(ctx, def.CodeErr, err.Error(), err)
		return
	}
	usr, err := ser
}
