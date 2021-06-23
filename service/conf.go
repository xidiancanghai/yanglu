package service

import "yanglu/service/model"

type ConfService struct {
}

func NewConfService() *ConfService {
	return &ConfService{}
}

func (cs *ConfService) GetConf() map[string]interface{} {
	return map[string]interface{}{
		"add_host":             model.AuthorityAddHost,
		"check_soft":           model.AuthorityCheckSoft,
		"create_security_task": model.AuthorityCreateSecurityTask,
		"add_user":             model.AuthorityAddNewUser,
		"check_log":            model.AuthorityCheckLog,
		"delete_user":          model.AuthorityDeleteUser,
		"create_smart_task":    model.AuthorityCreateSmartTask,
		"update_soft":          model.AuthorityUpdateSoft,
		"create_user_group":    model.AuthorityCreateUserGroup,
		"check_docker":         model.AuthorityCheckDocker,
	}
}
