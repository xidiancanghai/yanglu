package model

import (
	"errors"
	"time"
	"yanglu/service/data"

	"github.com/sirupsen/logrus"
)

type TaskItem struct {
	Id         int
	Ip         string
	TaskId     int
	Status     int
	UpdateTime int64
	CreateTime int64
}

func NewTaskItem() *TaskItem {
	return &TaskItem{}
}

func (t *TaskItem) TableName() string {
	return "task_item_info"
}

func (te *TaskItem) Create() error {
	if te.Ip == "" || te.TaskId == 0 {
		return errors.New("参数错误")
	}
	te.CreateTime = time.Now().Unix()
	te.UpdateTime = te.CreateTime
	tx := data.GetDB().Create(te)
	if tx.Error != nil {
		logrus.Error("Create err ", tx)
	}
	return tx.Error
}

func (t *TaskItem) BatchCreate(list []*TaskItem) error {
	if len(list) == 0 {
		return errors.New("列表长度为空")
	}
	db := data.GetDB()
	tx := db.Create(list)
	if tx.Error != nil {
		logrus.Error("BatchCreate error", tx)
	}
	return tx.Error
}

func (t *TaskItem) Updates(m map[string]interface{}) error {
	if len(m) == 0 || t.Id == 0 {
		return errors.New("参数错误")
	}
	m["update_time"] = time.Now().Unix()
	tx := data.GetDB().Model(t).Where(" id = ?", t.Id).Updates(m)
	if tx.Error != nil {
		logrus.Error("Updates err ", tx)
	}
	return tx.Error
}

func (t *TaskItem) GetTaskProgress(taskId int) (map[int]int, error) {
	sqll := "select status, count(distinct(ip)) as num from task_item_info where task_id = ? group by status "
	rows, err := data.GetDB().Raw(sqll, taskId).Rows()
	if err != nil {
		logrus.Error("GetTaskProgress err ", err)
		return nil, err
	}
	defer rows.Close()
	res := map[int]int{}
	for rows.Next() {
		var status int
		var num int
		rows.Scan(&status, &num)
		res[status] = num
	}
	return res, nil
}
