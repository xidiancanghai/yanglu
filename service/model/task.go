package model

import (
	"errors"
	"time"
	"yanglu/service/data"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const (
	TaskTypeFastTask     = 0
	TaskTypeTimedTask    = 1
	TaskTypeRepeatedTask = 2
)

type Task struct {
	Id          int `gorm:"primaryKey"`
	Name        string
	Ip          string
	Type        int
	Status      int // 0 未开始，1执行中，2执行结束
	ExecuTime   int64
	IntervalDay int
	UpdateTime  int64
	CreateTime  int64
}

func NewTask() *Task {
	return &Task{}
}

func (t *Task) TableName() string {
	return "task_info"
}

func (t *Task) Create() error {
	if t.Ip == "" {
		return errors.New("主机ip错误")
	}
	t.CreateTime = time.Now().Unix()
	t.UpdateTime = t.CreateTime
	tx := data.GetDB().Create(t)
	if tx.Error != nil {
		logrus.Error("Create err ", tx)
	}
	return tx.Error
}

func (t *Task) BatchCreate(list []*Task) error {
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

func (t *Task) Updates(m map[string]interface{}) error {
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

func (t *Task) GetTask(id int) (*Task, error) {
	if id == 0 {
		return nil, errors.New("参数错误")
	}
	rt := new(Task)
	tx := data.GetDB().Where(" id = ?", id).First(rt)
	if tx.Error != nil && tx.Error != gorm.ErrRecordNotFound {
		logrus.Error("GetTask err ", tx)
		return nil, tx.Error
	}
	return rt, nil
}

func (t *Task) GetSmartTask() ([]*Task, error) {
	list := []*Task{}

	sqll := " select * from task_info where execu_time >= unix_timestamp(now()) - 30 and execu_time <=  unix_timestamp(now()) + 30 and type in (1,2)"

	tx := data.GetDB().Model(t).Raw(sqll).Find(&list)
	if tx.Error != nil && tx.Error != gorm.ErrRecordNotFound {
		logrus.Error("GetTask err ", tx)
		return nil, tx.Error
	}
	return list, nil
}

func (t *Task) GetDetail() ([]*Task, error) {
	list := []*Task{}

	sqll := " select * from task_info where type in (1,2)"

	tx := data.GetDB().Model(t).Raw(sqll).Find(&list)
	if tx.Error != nil && tx.Error != gorm.ErrRecordNotFound {
		logrus.Error("GetTask err ", tx)
		return nil, tx.Error
	}
	return list, nil
}

func (t *Task) DeleteTask(ips []string) error {

	tx := data.GetDB().Model(t).Where(" ip in (?) ", ips).Delete(Task{})

	if tx.Error != nil {
		logrus.Error("DeleteTask err ", tx)
	}

	return tx.Error
}
