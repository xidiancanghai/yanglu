package model

import (
	"errors"
	"time"
	"yanglu/service/data"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type HostInfo struct {
	Id           int    `json:"-"`
	Ip           string `json:"ip"`
	Port         int    `json:"port"`
	SshUser      string `json:"ssh_user"`
	SshPasswd    string `json:"-"`
	Department   string `json:"department"`
	BusinessName string `json:"business_name"`
	SystemOs     string `json:"system_os"`
	UpdateTime   int64  `json:"-"`
	CreateTime   int64  `json:"-"`
}

func NewHostInfo() *HostInfo {
	return &HostInfo{}
}

func (h *HostInfo) TableName() string {
	return "host_info"
}

func (h *HostInfo) Create() error {
	if h.Ip == "" || h.Port == 0 || h.SshUser == "" || h.SshPasswd == "" {
		return errors.New("参数错误")
	}
	if h.CreateTime == 0 {
		h.CreateTime = time.Now().Unix()
		h.UpdateTime = h.CreateTime
	}
	db := data.GetDB()
	tx := db.Create(h)
	if tx.Error != nil {
		logrus.Error("insert error", tx)
	}
	return tx.Error
}

func (h *HostInfo) BatchCreate(list []*HostInfo) error {
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

func (h *HostInfo) Updates(m map[string]interface{}) error {
	if h.Id == 0 {
		return errors.New("参数错误")
	}
	tx := data.GetDB().Model(h).Where(" id = ?", h.Id).Updates(m)
	if tx.Error != nil {
		logrus.Error("Updates err", tx)
		return tx.Error
	}
	return nil
}

func (h *HostInfo) GetHostInfoById(id int) (*HostInfo, error) {
	if id == 0 {
		return nil, errors.New("参数错误")
	}
	rh := new(HostInfo)
	tx := data.GetDB().Where(" id = ? ", id).First(rh)
	if tx.Error != nil && tx.Error != gorm.ErrRecordNotFound {
		logrus.Error("GetHostInfoById err", tx)
		return nil, tx.Error
	}
	return rh, nil
}

func (h *HostInfo) GetHostInfoByIp(ip string) (*HostInfo, error) {
	if ip == "" {
		return nil, errors.New("参数错误")
	}
	rh := new(HostInfo)
	tx := data.GetDB().Where(" ip = ? ", ip).First(rh)
	if tx.Error != nil && tx.Error != gorm.ErrRecordNotFound {
		logrus.Error("GetHostInfoById err", tx)
		return nil, tx.Error
	}
	return rh, nil
}

func (h *HostInfo) GetHostsByNetworkSegment(ip string) ([]*HostInfo, error) {
	if ip == "" {
		return nil, errors.New("参数错误")
	}
	list := []*HostInfo{}
	sqll := "select * from host_info where ip like '" + ip + "%'"
	tx := data.GetDB().Model(h).Raw(sqll, ip).Find(&list)
	if tx.Error != nil && tx.Error != gorm.ErrRecordNotFound {
		logrus.Error("GetHostsByNetworkSegment err", tx)
		return nil, tx.Error
	}
	return list, nil
}

func (h *HostInfo) GetHostInfoByDepartment(department string) ([]*HostInfo, error) {
	if department == "" {
		return nil, errors.New("参数错误")
	}
	list := []*HostInfo{}
	tx := data.GetDB().Where(" department = ? ", department).Find(&list)
	if tx.Error != nil && tx.Error != gorm.ErrRecordNotFound {
		logrus.Error("GetHostInfoById err", tx)
		return nil, tx.Error
	}
	return list, nil
}

func (h *HostInfo) ListAll() ([]*HostInfo, error) {
	list := []*HostInfo{}
	tx := data.GetDB().Model(h).Find(&list)
	if tx.Error != nil && tx.Error != gorm.ErrRecordNotFound {
		logrus.Error("GetHostInfoById err", tx)
		return nil, tx.Error
	}
	return list, nil
}

func (h *HostInfo) SystemOsDistribute() ([]map[string]interface{}, error) {
	list := []map[string]interface{}{}
	sqll := "select system_os, count(*) as num from host_info group by system_os order by num desc"
	rows, err := data.GetDB().Raw(sqll).Rows()
	if err != nil && err != gorm.ErrRecordNotFound {
		logrus.Error("SystemOsDistribute err ", err)
		return list, err
	}
	if err != nil {
		return list, nil
	}
	defer rows.Close()
	for rows.Next() {
		var systemOs string
		var num int
		rows.Scan(&systemOs, &num)
		list = append(list, map[string]interface{}{
			"system_os": systemOs,
			"num":       num,
		})
	}
	return list, nil
}
