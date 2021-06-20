package service

import (
	"fmt"
	"yanglu/service/model"
)

type ActionLogService struct {
	u         *model.User
	actionLog *model.ActionLog
}

func NewActionLogService(uid int) *ActionLogService {
	u, _ := model.NewUser().GetUserById(uid)
	return &ActionLogService{
		u:         u,
		actionLog: model.NewActionLog(),
	}
}

func (as *ActionLogService) AddUser(newU *model.User) {
	as.actionLog.Uid = as.u.Uid
	as.actionLog.Detail = fmt.Sprintf("%s添加了用户%s", as.u.Name, newU.Name)
	as.actionLog.Create()
}

func (as *ActionLogService) Login() {
	as.actionLog.Uid = as.u.Uid
	as.actionLog.Detail = fmt.Sprintf("%s用户登陆", as.u.Name)
	as.actionLog.Create()
}

func (as *ActionLogService) AddHost(host *model.HostInfo) {
	as.actionLog.Uid = as.u.Uid
	as.actionLog.Detail = fmt.Sprintf("%s用户添加了主机%s", as.u.Name, host.Ip)
	as.actionLog.Create()
}

func (as *ActionLogService) BatchAddHost() {
	as.actionLog.Uid = as.u.Uid
	as.actionLog.Detail = fmt.Sprintf("%s用户批量添加了主机", as.u.Name)
	as.actionLog.Create()
}

func (as *ActionLogService) AddFastTask(task *model.Task) {
	as.actionLog.Uid = as.u.Uid
	as.actionLog.Detail = fmt.Sprintf("%s用户添加了安检任务%s", as.u.Name, task.Ip)
	as.actionLog.Create()
}

func (as *ActionLogService) AddTimedTask(task *model.Task) {
	as.actionLog.Uid = as.u.Uid
	as.actionLog.Detail = fmt.Sprintf("%s用户添加了定时安检任务%s", as.u.Name, task.Ip)
	as.actionLog.Create()
}

func (as *ActionLogService) AddRepeatedTask(task *model.Task) {
	as.actionLog.Uid = as.u.Uid
	as.actionLog.Detail = fmt.Sprintf("%s用户添加了重复安检任务%s", as.u.Name, task.Ip)
	as.actionLog.Create()
}

func (as *ActionLogService) List(lastId int) ([]*model.ActionLog, error) {
	return model.NewActionLog().List(lastId)
}
