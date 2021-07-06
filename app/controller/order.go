package controller

import (
	"yanglu/config"
	"yanglu/helper"

	"github.com/gin-gonic/gin"
)

type Order struct {
}

func NewOrder() *Order {
	return &Order{}
}

func (oc *Order) GetConfig(ctx *gin.Context) {
	helper.OKRsp(ctx, config.GetGoodsConfig())
}
