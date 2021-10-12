package controller

import (
	"errors"
	"yanglu/def"
	"yanglu/helper"
	"yanglu/service"
	"yanglu/service/model"

	"github.com/gin-gonic/gin"
)

type Log struct {
}

func NewLog() *Log {
	return &Log{}
}

func (l *Log) List(ctx *gin.Context) {

	params := &struct {
		LastId int `form:"last_id"  binding:"required"`
	}{}

	if err := ctx.ShouldBind(params); err != nil {
		helper.ErrRsp(ctx, def.CodeErr, "参数不正确", err)
		return
	}
	// 解析
	uid := ctx.GetInt("uid")
	// 先校验权限
	if !service.NewUserService().HasAuthority(uid, model.AuthorityCheckLog) {
		helper.ErrRsp(ctx, def.CodeErr, "你没有权限查看日志", errors.New("你没有权限查看日志"))
		return
	}
	//
	list, err := service.NewActionLogService(uid).List(params.LastId)
	if err != nil {
		helper.ErrRsp(ctx, def.CodeErr, err.Error(), err)
		return
	}
	helper.OKRsp(ctx, gin.H{
		"list": list,
	})
}

func (l *Log) SearchLog(ctx *gin.Context) {

	params := &struct {
		StartTime int    `form:"start_time"`
		EndTime   int    `form:"end_time"`
		Type      int    `form:"type"`
		Ip        string `form:"ip"`
	}{}

	if err := ctx.ShouldBind(params); err != nil {
		helper.ErrRsp(ctx, def.CodeErr, "参数不正确", err)
		return
	}
	// 解析
	uid := ctx.GetInt("uid")
	// 先校验权限
	if !service.NewUserService().HasAuthority(uid, model.AuthorityCheckLog) {
		helper.ErrRsp(ctx, def.CodeErr, "你没有权限查看日志", errors.New("你没有权限查看日志"))
		return
	}

	list, err := service.NewActionLogService(uid).SearchLog(params.StartTime, params.EndTime, params.Type, params.Ip)
	if err != nil {
		helper.ErrRsp(ctx, def.CodeErr, err.Error(), err)
		return
	}
	helper.OKRsp(ctx, gin.H{
		"list": list,
	})
}
