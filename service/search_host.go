package service

import (
	"errors"
	"yanglu/config"
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

func (s *SearchHostByIp) Search(uid int, ip string) ([]*model.HostInfo, error) {
	if ip == "" {
		return nil, errors.New("ip地址为空")
	}
	hs, err := model.NewHostInfo().GetHostsByNetworkSegment(ip)
	if err != nil {
		return nil, err
	}
	res := []*model.HostInfo{}
	if config.IsCloud() {
		for k, v := range hs {
			if v.Uid == uid {
				res = append(res, hs[k])
			}
		}
	} else {
		res = hs
	}
	return hs, nil
}

type SearchHostByDepartment struct {
}

func (s *SearchHostByDepartment) Search(uid int, department string) ([]*model.HostInfo, error) {
	if department == "" {
		return nil, errors.New("部门为空")
	}
	list, err := model.NewHostInfo().GetHostInfoByDepartment(department)
	if err != nil {
		return nil, err
	}
	res := []*model.HostInfo{}
	if config.IsCloud() {
		for k, v := range list {
			if v.Uid == uid {
				res = append(res, list[k])
			}
		}
	} else {
		res = list
	}
	return res, nil
}
