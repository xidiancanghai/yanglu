package controller

import (
	"errors"
	"yanglu/def"
	"yanglu/helper"
	"yanglu/service"
	"yanglu/service/model"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Task struct {
}

func NewTask() *Task {
	return &Task{}
}

func (t *Task) StartFastTask(ctx *gin.Context) {
	params := &struct {
		Ip string `form:"ip" binding:"required"`
	}{}

	if err := ctx.ShouldBind(params); err != nil {
		helper.ErrRsp(ctx, def.CodeErr, "参数不正确", err)
		return
	}
	uid := ctx.GetInt("uid")
	if !service.NewUserService().HasAuthority(uid, model.AuthorityCreateSecurityTask) {
		helper.ErrRsp(ctx, def.CodeErr, "你没有权限创建安全检查任务", errors.New("你没有权限创建安全检查任务"))
		return
	}

	// 添加机器
	ts := service.NewTaskService()
	task, err := ts.AddFastTask(params.Ip)
	if err != nil {
		helper.ErrRsp(ctx, def.CodeErr, err.Error(), err)
		return
	}
	go func() {
		err := ts.ExecuteTask(task)
		if err != nil {
			logrus.Error("StartFastTask err ", err)
		}
	}()
	service.NewActionLogService(uid).AddFastTask(task)
	// 执行任务
	helper.OKRsp(ctx, gin.H{
		"task_id": task.Id,
	})
}

func (t *Task) GetProgress(ctx *gin.Context) {
	params := &struct {
		TaskId int `form:"task_id" binding:"required"`
	}{}

	if err := ctx.ShouldBind(params); err != nil {
		helper.ErrRsp(ctx, def.CodeErr, "参数不正确", err)
		return
	}
	res, err := service.NewTaskService().GetProgress(params.TaskId)
	if err != nil {
		helper.ErrRsp(ctx, def.CodeErr, err.Error(), err)
		return
	}
	helper.OKRsp(ctx, res)
}

func (t *Task) AddTimedTask(ctx *gin.Context) {
	params := &struct {
		Ip        string `form:"ip" binding:"required"`
		ExecuTime int64  `form:"execu_time" binding:"required"`
		Name      string `form:"name" binding:"required"`
	}{}

	if err := ctx.ShouldBind(params); err != nil {
		helper.ErrRsp(ctx, def.CodeErr, "参数不正确", err)
		return
	}

	uid := ctx.GetInt("uid")
	if !service.NewUserService().HasAuthority(uid, model.AuthorityCreateSmartTask) {
		helper.ErrRsp(ctx, def.CodeErr, "你没有权限创建定期检查任务", errors.New("你没有权限创建定期检查任务"))
		return
	}

	task, err := service.NewTaskService().AddTimedTask(params.Ip, params.ExecuTime, params.Name)
	if err != nil {
		helper.ErrRsp(ctx, def.CodeErr, err.Error(), err)
		return
	}
	service.NewActionLogService(uid).AddTimedTask(task)
	// 执行任务
	helper.OKRsp(ctx, gin.H{
		"task_id": task.Id,
	})
}

func (t *Task) AddRepeatTask(ctx *gin.Context) {
	params := &struct {
		Ip        string `form:"ip" binding:"required"`
		ExecuTime int64  `form:"execu_time" binding:"required"`
		Interval  int    `form:"interval" binding:"required"`
		Name      string `form:"name" binding:"required"`
	}{}

	if err := ctx.ShouldBind(params); err != nil {
		helper.ErrRsp(ctx, def.CodeErr, "参数不正确", err)
		return
	}

	uid := ctx.GetInt("uid")
	if !service.NewUserService().HasAuthority(uid, model.AuthorityCreateSmartTask) {
		helper.ErrRsp(ctx, def.CodeErr, "你没有权限创建定期检查任务", errors.New("你没有权限创建定期检查任务"))
		return
	}

	task, err := service.NewTaskService().AddRepeatTask(params.Ip, params.ExecuTime, params.Interval, params.Name)
	if err != nil {
		helper.ErrRsp(ctx, def.CodeErr, err.Error(), err)
		return
	}
	service.NewActionLogService(uid).AddRepeatedTask(task)
	// 执行任务
	helper.OKRsp(ctx, gin.H{
		"task_id": task.Id,
	})
}

func (t *Task) GetDetail(ctx *gin.Context) {
	list, err := service.NewTaskService().GetDetail()
	if err != nil {
		helper.ErrRsp(ctx, def.CodeErr, err.Error(), err)
		return
	}
	helper.OKRsp(ctx, list)
}

func (t *Task) CheckInfo(ctx *gin.Context) {
	res, err := service.NewTaskService().GetHostCheckStatus()
	if err != nil {
		helper.ErrRsp(ctx, def.CodeErr, err.Error(), err)
		return
	}
	helper.OKRsp(ctx, res)
}
