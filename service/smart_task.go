package service

import (
	"fmt"
	"time"
	"yanglu/service/model"

	"github.com/sirupsen/logrus"
)

type SmartTask struct {
}

func NewSmartTask() *SmartTask {
	return &SmartTask{}
}

func (st *SmartTask) TaskCheck() {

	list, _ := model.NewTask().GetSmartTask()
	if len(list) == 0 {
		return
	}
	t := time.Now()
	log := logrus.WithField("start = ", fmt.Sprintf("%d:%d:%d", t.Hour(), t.Minute(), t.Second()))
	log.Info("start-----")
	for k, v := range list {
		log.Info(" task = ", *list[k])
		if v.Type == model.TaskTypeTimedTask {
			diff := v.ExecuTime - t.Unix()
			// 执行任务
			if diff > 0 && diff < 30 {
				NewTaskService().ExecuteTask(list[k])
			}
		} else if v.Type == model.TaskTypeRepeatedTask {
			diff := v.ExecuTime - t.Unix()
			// 执行任务
			if diff > 0 && diff < 30 {
				NewTaskService().ExecuteTask(list[k])
				v.ExecuTime = v.ExecuTime + int64(v.IntervalDay)*3600*24
				v.Updates(map[string]interface{}{
					"execu_time": v.ExecuTime,
				})
			}

		}
	}
	t = time.Now()
	log.Info("end--", fmt.Sprintf("%d:%d:%d", t.Hour(), t.Minute(), t.Second()))
}
