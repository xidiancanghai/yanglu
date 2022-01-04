package controller

import (
	"errors"
	"time"
	"yanglu/def"
	"yanglu/helper"
	"yanglu/service"
	"yanglu/service/model"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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
		StartTime string `form:"start_time"`
		EndTime   string `form:"end_time"`
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

	startTime := 0
	if params.StartTime != "" {
		temp, err := time.ParseInLocation("2006-01-02", params.StartTime, time.Local)
		logrus.Info("err = ", err)
		startTime = int(temp.Unix())
	}
	endTime := 0
	if params.EndTime != "" {
		temp, _ := time.ParseInLocation("2006-01-02", params.EndTime, time.Local)
		endTime = int(temp.Unix())
	}

	logrus.Info(startTime, " ", endTime)

	list, err := service.NewActionLogService(uid).SearchLog(startTime, endTime, params.Type, params.Ip)
	if err != nil {
		helper.ErrRsp(ctx, def.CodeErr, err.Error(), err)
		return
	}
	helper.OKRsp(ctx, gin.H{
		"list": list,
	})
}
