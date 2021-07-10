package controller

import (
	"yanglu/config"
	"yanglu/def"
	"yanglu/helper"
	"yanglu/service"

	"github.com/gin-gonic/gin"
)

type UtilController struct {
}

func NewUtilController() *UtilController {
	return &UtilController{}
}

func (uc *UtilController) GetCaptchaId(ctx *gin.Context) {
	helper.OKRsp(ctx, gin.H{
		"id": service.NewUtilService().GetCaptchaId(),
	})
}

func (uc *UtilController) GetCaptcha(ctx *gin.Context) {

	params := &struct {
		Id string `form:"id"  binding:"required"`
	}{}

	if err := ctx.ShouldBind(params); err != nil {
		helper.ErrRsp(ctx, def.CodeErr, "参数不正确", err)
		return
	}

	buf, err := service.NewUtilService().GetCaptcha(params.Id)
	if err != nil {
		helper.ErrRsp(ctx, def.CodeErr, err.Error(), err)
		return
	}
	ctx.Writer.Header().Add("Content-Type", "image/png")
	_, err = buf.WriteTo(ctx.Writer)
	if err != nil {
		helper.ErrRsp(ctx, def.CodeErr, err.Error(), err)
		return
	}
}

func (uc *UtilController) GetSystemInfo(ctx *gin.Context) {
	helper.OKRsp(ctx, gin.H{
		"max_node": config.LicenseInfoConf.NodeMax,
		"edition":  config.LicenseInfoConf.Edition,
	})
}

func (uc *UtilController) UploadImages(ctx *gin.Context) {
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		return
	}
	if err != nil {
		return
	}
	defer file.Close()

	fileName, err := service.NewUtilService().UploadImage(file, header.Filename)
	if err != nil {
		helper.ErrRsp(ctx, def.CodeErr, err.Error(), err)
		return
	}
	helper.OKRsp(ctx, gin.H{"file_name": fileName})
}

func (uc *UtilController) DownloadImage(ctx *gin.Context) {

	params := &struct {
		FileName string `form:"file_name"  binding:"required"`
	}{}

	if err := ctx.ShouldBind(params); err != nil {
		helper.ErrRsp(ctx, def.CodeErr, "参数不正确", err)
		return
	}

	err := service.NewUtilService().DownloadImage(params.FileName, ctx)
	if err != nil {
		helper.ErrRsp(ctx, def.CodeErr, err.Error(), err)
		return
	}
}
