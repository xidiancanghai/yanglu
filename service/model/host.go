package model

import (
	"errors"
	"fmt"
	"time"
	"yanglu/config"
	"yanglu/service/data"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type HostInfo struct {
	Id           int    `json:"-" gorm:"primaryKey"`
	Ip           string `json:"ip"`
	Port         int    `json:"port"`
	SshUser      string `json:"ssh_user"`
	SshPasswd    string `json:"-"`
	Department   string `json:"department"`
	BusinessName string `json:"business_name"`
	SystemOs     string `json:"system_os"`
	Uid          int    `json:"-"`
	IsDelete     int    `json:"-"`
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
	if h.Id == 0 || len(m) == 0 {
		return errors.New("参数错误")
	}
	m["update_time"] = time.Now().Unix()
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
	tx := data.GetDB().Where(" ip = ? and is_delete = 0 ", ip).First(rh)
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
	sqll := "select * from host_info where ip like '" + ip + "%' and is_delete = 0 "
	tx := data.GetDB().Model(h).Raw(sqll, ip).Find(&list)
	if tx.Error != nil && tx.Error != gorm.ErrRecordNotFound {
		logrus.Error("GetHostsByNetworkSegment err", tx)
		return nil, tx.Error
	}
	return list, nil
}

func (h *HostInfo) GetHostsByIps(ips []string) ([]*HostInfo, error) {
	if len(ips) == 0 {
		return nil, errors.New("参数错误")
	}
	list := []*HostInfo{}
	sqll := "select * from host_info where ip in (?) and is_delete = 0 "
	tx := data.GetDB().Model(h).Raw(sqll, ips).Find(&list)
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
	tx := data.GetDB().Where(" department = ? and is_delete = 0", department).Find(&list)
	if tx.Error != nil && tx.Error != gorm.ErrRecordNotFound {
		logrus.Error("GetHostInfoById err", tx)
		return nil, tx.Error
	}
	return list, nil
}

func (h *HostInfo) ListAll() ([]*HostInfo, error) {
	list := []*HostInfo{}
	tx := data.GetDB().Model(h).Where(" is_delete = 0 ").Find(&list)
	if tx.Error != nil && tx.Error != gorm.ErrRecordNotFound {
		logrus.Error("GetHostInfoById err", tx)
		return nil, tx.Error
	}
	return list, nil
}

func (h *HostInfo) SystemOsDistribute(uid int) ([]map[string]interface{}, error) {
	list := []map[string]interface{}{}
	sqll := "select system_os, count(*) as num from host_info where is_delete = 0 group by system_os order by num desc"
	if config.IsCloud() {
		sqll = fmt.Sprintf("select system_os, count(*) as num from host_info where uid = %d and is_delete = 0 group by system_os order by num desc", uid)
	}
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

func (h *HostInfo) BatchDelete(ips []string) error {

	sqll := "update " + h.TableName() + " set is_delete = 1, update_time = unix_timestamp(now()) where ip in (?)"

	_, err := data.GetDB().Model(h).Raw(sqll, ips).Rows()
	if err != nil {
		logrus.Error("BatchDelete err", err)
		return err
	}
	return nil
}

func (h *HostInfo) Delete(ips []string) error {

	tx := data.GetDB().Model(h).Where(" ip in (?) ", ips).Delete(HostInfo{})

	if tx.Error != nil {
		logrus.Error("Delete err ", tx)
	}

	return tx.Error
}
