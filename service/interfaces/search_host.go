package interfaces

import "yanglu/service/model"

type SearchHost interface {
	Search(uid int, condition string) ([]*model.HostInfo, error)
}
