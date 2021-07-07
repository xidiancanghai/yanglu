package service

import (
	"encoding/json"
	"errors"
	"strings"
	"sync"
	"time"
	"yanglu/config"
	"yanglu/service/model"

	"github.com/sirupsen/logrus"
)

type TaskService struct {
}

func NewTaskService() *TaskService {
	return &TaskService{}
}

func (ts *TaskService) GetHostInfoByIp(ip string) ([]*model.HostInfo, error) {
	ss := strings.Split(ip, ".")
	list := []*model.HostInfo{}
	var err error = nil
	// 网段
	if len(ss) == 3 {
		list, err = model.NewHostInfo().GetHostsByNetworkSegment(ip)
		if err != nil {
			logrus.Error("GetHostInfoByIp err ", err)
			return list, err
		}
	} else if len(ss) == 4 {
		host, err := model.NewHostInfo().GetHostInfoByIp(ip)
		if err != nil {
			logrus.Error("GetHostInfoByIp err ", err)
			return list, err
		}
		list = []*model.HostInfo{host}
	} else {
		return list, errors.New("错误的ip")
	}
	return list, nil
}

func (ts *TaskService) AddFastTask(ip string) (*model.Task, error) {
	task := model.Task{
		Ip:        ip,
		Type:      model.TaskTypeFastTask,
		Status:    0,
		ExecuTime: time.Now().Unix(),
	}
	err := task.Create()
	if err != nil {
		logrus.Error("AddTimedTask err ", err)
		return nil, err
	}
	return &task, nil
}

func (ts *TaskService) AddTimedTask(ip string, execuTime int64, name string) (*model.Task, error) {
	if config.LicenseInfoConf.SmartTask == 0 {
		return nil, errors.New("当前系统没有添加智能任务权限")
	}
	task := model.Task{
		Ip:        ip,
		Type:      model.TaskTypeTimedTask,
		Status:    0,
		ExecuTime: execuTime,
		Name:      name,
	}
	err := task.Create()
	if err != nil {
		logrus.Error("AddTimedTask err ", err)
		return nil, err
	}
	return &task, nil
}

func (ts *TaskService) AddRepeatTask(ip string, execuTime int64, interval int, name string) (*model.Task, error) {
	if config.LicenseInfoConf.SmartTask == 0 {
		return nil, errors.New("当前系统没有添加智能任务权限")
	}
	task := model.Task{
		Ip:          ip,
		Type:        model.TaskTypeRepeatedTask,
		Status:      0,
		ExecuTime:   execuTime,
		IntervalDay: interval,
		Name:        name,
	}
	err := task.Create()
	if err != nil {
		logrus.Error("AddRepeatTask err ", err)
		return nil, err
	}
	return &task, nil
}

func (ts *TaskService) ExecuteTask(task *model.Task) error {

	if config.LicenseInfoConf.SmartTask == 0 {
		return errors.New("当前系统没有添加智能任务权限")
	}

	if config.LicenseInfoConf.NodeMax < NewHostInfoService().GetHostNum() {
		return errors.New("当前系统机器数量超过了最大限制")
	}

	// 首先找出所有的ip
	if task.Ip == "" {
		return errors.New("ip错误")
	}
	hostList, err := ts.GetHostInfoByIp(task.Ip)
	if len(hostList) == 0 {
		return errors.New("找不到机器")
	}

	taskItems := make([]*model.TaskItem, len(hostList))
	t := time.Now().Unix()
	for i := 0; i < len(hostList); i++ {
		taskItems[i] = &model.TaskItem{
			Ip:         hostList[i].Ip,
			TaskId:     task.Id,
			Status:     0,
			CreateTime: t,
			UpdateTime: t,
		}
	}

	err = model.NewTaskItem().BatchCreate(taskItems)
	if err != nil {
		logrus.Error("Execute err ", err)
	}

	// 任务修改未执行中
	err = task.Updates(map[string]interface{}{
		"status": 1,
	})

	if err != nil {
		logrus.Error("Execute err ", err)
		return err
	}

	var wg sync.WaitGroup
	go func() {
		for k, v := range taskItems {
			wg.Add(1)
			go func() {
				defer wg.Done()
				err := ts.StartTask(taskItems[k])
				if err != nil {
					logrus.Error("Execute err ", err)
				}
			}()
			// 正在进行中
			v.Status = 1
			err = v.Updates(map[string]interface{}{
				"status": 1,
			})
			if err != nil {
				logrus.Error("Execute err ", err)
				continue
			}
		}
	}()
	wg.Wait()

	err = task.Updates(map[string]interface{}{
		"status": 2,
	})

	if err != nil {
		logrus.Error("Execute err ", err)
		return err
	}

	return nil
}

func (ts *TaskService) EnvironmentOk(hostInfo *model.HostInfo) (bool, error) {
	hostService := NewHostInfoService()
	// 先检测文件是否都存在
	cmds := []string{
		`
	cd /var/trivy;file="/var/trivy/trivy";if [ ! -f "$file" ];then  echo 1; else  echo 0; fi
	`,
		`
	cd /var/trivy;file="/root/.cache/trivy/db/metadata.json";if [ ! -f "$file" ];then  echo 1; else  echo 0; fi
	`,
		`
	cd /var/trivy;file="/root/.cache/trivy/db/trivy.db";if [ ! -f "$file" ];then  echo 1; else  echo 0; fi
	`,
	}

	for _, cmd := range cmds {
		res, err := hostService.Cmd(hostInfo, cmd)
		if err != nil {
			logrus.Error("StartTask err ", err)
			return false, err
		}
		if strings.TrimSpace(res) == "1" {
			return false, nil
		}
	}

	return true, nil
}

func (ts *TaskService) StartTask(taskItem *model.TaskItem) error {
	hostInfo, _ := model.NewHostInfo().GetHostInfoByIp(taskItem.Ip)

	ok, err := ts.EnvironmentOk(hostInfo)

	if err != nil {
		logrus.Error("StartTask err ", err)
		return err
	}

	hostService := NewHostInfoService()
	if !ok {
		logrus.WithFields(logrus.Fields{
			"host": *hostInfo,
		}).Info("StartTask")
		hostService.Prepare([]*model.HostInfo{hostInfo})

		count := 0
		for {
			ok, _ = ts.EnvironmentOk(hostInfo)
			if ok {
				break
			}
			count++
			if count > 20 {
				break
			}
			time.Sleep(20 * time.Second)
		}
		if !ok {
			return errors.New("环境准备失败，请检查网络")
		}
	}

	cmd := `
	cd /var/trivy; rm -rf results.json;./trivy fs --skip-update -f json -o results.json /
	`

	_, err = hostService.Cmd(hostInfo, cmd)

	if err != nil {
		logrus.Error("StartTask err ", err)
		return err
	}

	res, err := hostService.GetResult(hostInfo)

	if err != nil {
		return err
	}

	list := []model.TrivyResult{}

	json.Unmarshal([]byte(res), &list)

	if len(list) == 0 {
		return errors.New("找不到日志")
	}

	logList := []*model.VulnerabilityLog{}
	if err != nil {
		return err
	}
	t := time.Now().Unix()
	for _, v1 := range list {
		for _, v2 := range v1.Vulnerabilities {
			logList = append(logList, &model.VulnerabilityLog{
				Ip:               taskItem.Ip,
				TaskId:           taskItem.TaskId,
				TaskItemId:       taskItem.Id,
				VulnerabilityId:  v2.VulnerabilityId,
				PkgName:          v2.PkgName,
				InstalledVersion: v2.InstalledVersion,
				FixedVersion:     v2.FixedVersion,
				Severity:         v2.Severity,
				CreateTime:       t,
			})
		}
	}
	// 标记状态未已完成
	taskItem.Updates(map[string]interface{}{
		"status": 2,
	})
	if len(logList) == 0 {
		return errors.New("日志错误")
	}
	err = model.NewVulnerabilityLog().BatchCreate(logList)
	if err != nil {
		logrus.Error("StartTask err ", err)
	}

	return nil
}

func (ts *TaskService) GetProgress(taskId int) (interface{}, error) {
	info, err := model.NewTaskItem().GetTaskProgress(taskId)
	if err != nil {
		logrus.Error("GetTaskProgress err ", err)
		return nil, err
	}
	all := 0
	res := map[string]interface{}{}
	for k, v := range info {
		all += v
		if k == 0 {
			res["no_start"] = v
		} else if k == 1 {
			res["checking"] = v
		} else if k == 2 {
			res["finished"] = v
		}
	}
	res["all"] = all
	return res, nil
}

func (ts *TaskService) GetDetail() (interface{}, error) {
	list, err := model.NewTask().GetDetail()
	if err != nil {
		logrus.Error("GetDetail err ", err)
		return nil, err
	}
	res := struct {
		All         int `json:"all"`
		Checking    int `json:"checking"`
		QueueTask   int `json:"queue_task"`
		PlaningTask int `json:"planing_task"`
	}{}
	res.All = len(list)
	for _, v := range list {
		if v.Status == 1 {
			res.Checking++
		}
		if v.ExecuTime > time.Now().Unix() {
			res.QueueTask++
		}
		if v.Type == model.TaskTypeTimedTask {
			res.PlaningTask++
		}
	}
	return res, nil
}

func (ts *TaskService) GetHostCheckStatus() (interface{}, error) {
	hosts, err := model.NewHostInfo().ListAll()
	if err != nil {
		logrus.Error("ListAll err ", err)
		return nil, err
	}
	if len(hosts) == 0 {
		return nil, errors.New("暂未有任何机器")
	}
	ips := make([]string, len(hosts))
	for k, v := range hosts {
		ips[k] = v.Ip
	}
	checkStatus, err := model.NewVulnerabilityLog().HostCheckStatus(ips)
	if err != nil {
		logrus.Error("ListAll err ", err)
		return nil, err
	}
	res := map[string]int{}
	res["has_checked"] = len(checkStatus)
	res["wait_check"] = len(hosts) - len(checkStatus)
	return res, nil
}
