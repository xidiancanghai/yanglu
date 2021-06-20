package controller

import (
	"errors"
	"yanglu/def"
	"yanglu/helper"
	"yanglu/service"
	"yanglu/service/model"

	"github.com/gin-gonic/gin"
)

type Host struct {
}

func NewHost() *Host {
	return &Host{}
}

func (hc *Host) Add(ctx *gin.Context) {
	params := &struct {
		Ip         string `form:"ip" binding:"required"`
		Port       int    `form:"port" binding:"required"`
		SshUser    string `form:"ssh_user" binding:"required"`
		SshPasswd  string `form:"ssh_passwd" binding:"required"`
		Department string `form:"department"`
	}{}

	if err := ctx.ShouldBind(params); err != nil {
		helper.ErrRsp(ctx, def.CodeErr, "参数不正确", err)
		return
	}
	uid := ctx.GetInt("uid")

	// 先校验权限
	if !service.NewUserService().HasAuthority(uid, model.AuthorityAddHost) {
		helper.ErrRsp(ctx, def.CodeErr, "你没有权限添加主机", errors.New("你没有权限添加主机"))
		return
	}
	// 添加机器
	host, err := service.NewHostInfoService().Add(params.Ip, params.Port, params.SshUser, params.SshPasswd, params.Department)
	if err != nil {
		helper.ErrRsp(ctx, def.CodeErr, err.Error(), err)
		return
	}
	service.NewActionLogService(uid).AddHost(host)
	helper.OKRsp(ctx, gin.H{})
}

func (hc *Host) BatchAdd(ctx *gin.Context) {
	file, _, err := ctx.Request.FormFile("file")
	if err != nil {
		return
	}
	if err != nil {
		return
	}
	defer file.Close()

	uid := ctx.GetInt("uid")
	// 先校验权限
	if !service.NewUserService().HasAuthority(uid, model.AuthorityAddHost) {
		helper.ErrRsp(ctx, def.CodeErr, "你没有权限添加主机", errors.New("你没有权限添加主机"))
		return
	}

	excelService, err := service.NewExcelService(file)
	if err != nil {
		helper.ErrRsp(ctx, def.CodeErr, err.Error(), err)
		return
	}
	hostList, err := excelService.GetHostInfos()
	if err != nil {
		helper.ErrRsp(ctx, def.CodeErr, err.Error(), err)
		return
	}
	err = service.NewHostInfoService().BatchAdd(hostList)
	if err != nil {
		helper.ErrRsp(ctx, def.CodeErr, err.Error(), err)
		return
	}
	service.NewActionLogService(uid).BatchAddHost()
	helper.OKRsp(ctx, gin.H{})
}

func (hc *Host) UpdateDepartment(ctx *gin.Context) {
	params := &struct {
		Ip         string `form:"ip" binding:"required"`
		Department string `form:"department" binding:"required"`
	}{}

	if err := ctx.ShouldBind(params); err != nil {
		helper.ErrRsp(ctx, def.CodeErr, "参数不正确", err)
		return
	}
	uid := ctx.GetInt("uid")
	// 先校验权限
	if !service.NewUserService().HasAuthority(uid, model.AuthorityAddHost) {
		helper.ErrRsp(ctx, def.CodeErr, "你没有权限添加主机", errors.New("你没有权限添加主机"))
		return
	}
	err := service.NewHostInfoService().UpdateDepartment(params.Ip, params.Department)
	if err != nil {
		helper.ErrRsp(ctx, def.CodeErr, err.Error(), err)
		return
	}
	helper.OKRsp(ctx, gin.H{})
}

func (hc *Host) SearchHost(ctx *gin.Context) {
	params := &struct {
		Type      int    `form:"type"`
		Condition string `form:"condition" binding:"required"`
	}{}

	if err := ctx.ShouldBind(params); err != nil {
		helper.ErrRsp(ctx, def.CodeErr, "参数不正确", err)
		return
	}
	//uid := ctx.GetInt("uid")
	searchHost := service.NewSearchHostFactory().CreateSearch(params.Type)
	if searchHost == nil {
		helper.ErrRsp(ctx, def.CodeErr, "请求类型错误", nil)
		return
	}
	list, err := searchHost.Search(params.Condition)
	if err != nil {
		helper.ErrRsp(ctx, def.CodeErr, err.Error(), err)
		return
	}
	helper.OKRsp(ctx, gin.H{
		"list": list,
	})
}

func (hc *Host) GetVulnerabilityInfo(ctx *gin.Context) {
	params := &struct {
		Ip string `form:"ip" binding:"required"`
	}{}

	if err := ctx.ShouldBind(params); err != nil {
		helper.ErrRsp(ctx, def.CodeErr, "参数不正确", err)
		return
	}
	vs, err := service.NewVulnerabilityService(params.Ip)
	if err != nil {
		helper.ErrRsp(ctx, def.CodeErr, err.Error(), err)
		return
	}
	list, err := vs.GetGetVulnerabilityInfo()
	if err != nil {
		helper.ErrRsp(ctx, def.CodeErr, err.Error(), err)
		return
	}
	helper.OKRsp(ctx, gin.H{
		"list": list,
	})
}
