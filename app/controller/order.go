package controller

import (
	"fmt"
	"yanglu/config"
	"yanglu/def"
	"yanglu/helper"
	"yanglu/service"

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

func (oc *Order) Create(ctx *gin.Context) {
	params := &struct {
		NodeNum      int    `form:"node_num" binding:"required"`
		MonthNum     int    `form:"month_num" binding:"required"`
		DiscountCode string `form:"discount_code"`
	}{}

	if err := ctx.ShouldBind(params); err != nil {
		helper.ErrRsp(ctx, def.CodeErr, "参数不正确", err)
		return
	}
	uid := ctx.GetInt("uid")

	uid = 30

	data, err := service.NewOrderService(uid).Create(params.NodeNum, params.MonthNum, params.DiscountCode)

	if err != nil {
		helper.ErrRsp(ctx, def.CodeErr, err.Error(), err)
		return
	}

	ctx.Writer.Header().Set("Content-Type", "image/png")
	ctx.Writer.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
	ctx.Writer.Write(data)
}
