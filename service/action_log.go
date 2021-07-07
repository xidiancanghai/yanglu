package service

import (
	"fmt"
	"yanglu/config"
	"yanglu/service/model"
)

type ActionLogService struct {
	uid       int
	name      string
	actionLog *model.ActionLog
}

func NewActionLogService(uid int) *ActionLogService {
	name := ""
	if config.IsCloud() {
		u, _ := model.NewCloudUser().GetUser(map[string]interface{}{
			"uid": uid,
		})
		if u != nil {
			name = u.Phone
		}
	} else {
		u, _ := model.NewUser().GetUserById(uid)
		if u != nil {
			name = u.Name
		}
	}
	return &ActionLogService{
		uid:       uid,
		name:      name,
		actionLog: model.NewActionLog(),
	}
}

func (as *ActionLogService) AddUser(newU *model.User) {
	as.actionLog.Uid = as.uid
	as.actionLog.Detail = fmt.Sprintf("%s添加了用户%s", as.name, newU.Name)
	as.actionLog.Create()
}

func (as *ActionLogService) Login() {
	as.actionLog.Uid = as.uid
	as.actionLog.Detail = fmt.Sprintf("%s用户登陆", as.name)
	as.actionLog.Create()
}

func (as *ActionLogService) AddHost(host *model.HostInfo) {
	as.actionLog.Uid = as.uid
	as.actionLog.Detail = fmt.Sprintf("%s用户添加了主机%s", as.name, host.Ip)
	as.actionLog.Create()
}

func (as *ActionLogService) BatchAddHost() {
	as.actionLog.Uid = as.uid
	as.actionLog.Detail = fmt.Sprintf("%s用户批量添加了主机", as.name)
	as.actionLog.Create()
}

func (as *ActionLogService) AddFastTask(task *model.Task) {
	as.actionLog.Uid = as.uid
	as.actionLog.Detail = fmt.Sprintf("%s用户添加了安检任务%s", as.name, task.Ip)
	as.actionLog.Create()
}

func (as *ActionLogService) AddTimedTask(task *model.Task) {
	as.actionLog.Uid = as.uid
	as.actionLog.Detail = fmt.Sprintf("%s用户添加了定时安检任务%s", as.name, task.Ip)
	as.actionLog.Create()
}

func (as *ActionLogService) AddRepeatedTask(task *model.Task) {
	as.actionLog.Uid = as.uid
	as.actionLog.Detail = fmt.Sprintf("%s用户添加了重复安检任务%s", as.name, task.Ip)
	as.actionLog.Create()
}

func (as *ActionLogService) List(lastId int) (interface{}, error) {
	list, err := model.NewActionLog().List(lastId)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return list, nil
	}
	res := make([]struct {
		*model.ActionLog
		Name string `json:"name"`
	}, len(list))
	uids := make([]int, len(list))
	names := map[int]string{}

	for k, v := range list {
		uids[k] = v.Uid
	}

	if config.IsCloud() {
		names = model.NewCloudUser().GetUserNames(uids)
	} else {
		names = model.NewUser().GetUserNames(uids)
	}

	for k, v := range list {
		res[k].ActionLog = list[k]
		res[k].Name = names[v.Uid]
	}

	return res, nil
}
