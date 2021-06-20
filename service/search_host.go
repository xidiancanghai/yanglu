package service

import (
	"errors"
	"yanglu/service/interfaces"
	"yanglu/service/model"
)

type SearchHostFactory struct {
}

func NewSearchHostFactory() *SearchHostFactory {
	return &SearchHostFactory{}
}

func (sf *SearchHostFactory) CreateSearch(searchType int) interfaces.SearchHost {
	switch searchType {
	case 0:
		return &SearchHostByIp{}
	case 1:
		return &SearchHostByDepartment{}
	}
	return nil
}

type SearchHostByIp struct {
}

func (s *SearchHostByIp) Search(ip string) ([]*model.HostInfo, error) {
	if ip == "" {
		return nil, errors.New("ip地址为空")
	}
	hs, err := model.NewHostInfo().GetHostInfoByIp(ip)
	if err != nil {
		return nil, err
	}
	return []*model.HostInfo{hs}, nil
}

type SearchHostByDepartment struct {
}

func (s *SearchHostByDepartment) Search(department string) ([]*model.HostInfo, error) {
	if department == "" {
		return nil, errors.New("部门为空")
	}
	list, err := model.NewHostInfo().GetHostInfoByDepartment(department)
	if err != nil {
		return nil, err
	}
	return list, nil
}
