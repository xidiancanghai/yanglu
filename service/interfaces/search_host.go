package interfaces

import "yanglu/service/model"

type SearchHost interface {
	Search(condition string) ([]*model.HostInfo, error)
}
