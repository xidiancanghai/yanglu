package controller

import (
	"yanglu/helper"
	"yanglu/service"

	"github.com/gin-gonic/gin"
)

type ConfController struct {
}

func NewConfController() *ConfController {
	return &ConfController{}
}

func (c *ConfController) GetConf(ctx *gin.Context) {
	helper.OKRsp(ctx, service.NewConfService().GetConf())
}
