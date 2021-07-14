package service

import (
	"yanglu/config"
	"yanglu/service/model"
)

// type SearchHostFactory struct {
// }

// func NewSearchHostFactory() *SearchHostFactory {
// 	return &SearchHostFactory{}
// }

// func (sf *SearchHostFactory) CreateSearch(searchType int) interfaces.SearchHost {
// 	switch searchType {
// 	case 0:
// 		return &SearchHostByIp{}
// 	case 1:
// 		return &SearchHostByDepartment{}
// 	}
// 	return nil
// }

// type SearchHostByIp struct {
// }

// func (s *SearchHostByIp) Search(uid int, ip string) ([]*model.HostInfo, error) {
// 	if ip == "" {
// 		return nil, errors.New("ip地址为空")
// 	}
// 	hs, err := model.NewHostInfo().GetHostsByNetworkSegment(ip)
// 	if err != nil {
// 		return nil, err
// 	}
// 	res := []*model.HostInfo{}
// 	if config.IsCloud() {
// 		for k, v := range hs {
// 			if v.Uid == uid {
// 				res = append(res, hs[k])
// 			}
// 		}
// 	} else {
// 		res = hs
// 	}
// 	return hs, nil
// }

// type SearchHostByDepartment struct {
// }

// func (s *SearchHostByDepartment) Search(uid int, department string) ([]*model.HostInfo, error) {
// 	if department == "" {
// 		return nil, errors.New("部门为空")
// 	}
// 	list, err := model.NewHostInfo().GetHostInfoByDepartment(department)
// 	if err != nil {
// 		return nil, err
// 	}
// 	res := []*model.HostInfo{}
// 	if config.IsCloud() {
// 		for k, v := range list {
// 			if v.Uid == uid {
// 				res = append(res, list[k])
// 			}
// 		}
// 	} else {
// 		res = list
// 	}
// 	return res, nil
// }

type SearchHostService struct {
	uid       int
	condition string
	list      []*model.HostInfo
}

func NewSearchHostService(uid int, condition string) *SearchHostService {
	return &SearchHostService{
		uid:       uid,
		condition: condition,
		list:      make([]*model.HostInfo, 0),
	}
}

func (ss *SearchHostService) SearchHost() ([]*model.HostInfo, error) {

	// 搜IP
	list1, err := model.NewHostInfo().GetHostsByNetworkSegment(ss.condition)
	if err != nil {
		return nil, err
	}
	ss.list = append(ss.list, list1...)

	// 搜部门
	list2, err := model.NewHostInfo().GetHostInfoByDepartment(ss.condition)
	if err != nil {
		return nil, err
	}
	ss.list = append(ss.list, list2...)

	// 搜软件包
	ips, _ := model.NewVulnerabilityLog().ListHostBySoftName(ss.condition)
	list3, _ := model.NewHostInfo().GetHostsByIps(ips)
	ss.list = append(ss.list, list3...)

	// 去重复
	uniqIps := map[string]bool{}
	uniqList := []*model.HostInfo{}
	for _, v := range ss.list {
		uniqIps[v.Ip] = true
	}

	for k, v := range ss.list {
		if config.IsCloud() && v.Uid != ss.uid {
			continue
		}
		if uniqIps[v.Ip] {
			uniqList = append(uniqList, ss.list[k])
			delete(uniqIps, v.Ip)
		}
	}

	return uniqList, nil
}
