package controller

import (
	"yanglu/def"
	"yanglu/helper"
	"yanglu/service"

	"github.com/gin-gonic/gin"
)

type ArticleController struct {
}

func NewArticleController() *ArticleController {
	return &ArticleController{}
}

func (ac *ArticleController) Add(ctx *gin.Context) {
	params := &struct {
		Title   string   `json:"title" binding:"required"`
		Tag     string   `json:"tag" binding:"required"`
		Content string   `json:"content" binding:"required"`
		Photos  []string `json:"photos"`
	}{}

	if err := ctx.ShouldBind(params); err != nil {
		helper.ErrRsp(ctx, def.CodeErr, "参数不正确", err)
		return
	}
	uid := ctx.GetInt("uid")
	as := service.NewArticleService(uid)

	err := as.Add(params.Title, params.Tag, params.Content, params.Photos)

	if err != nil {
		helper.ErrRsp(ctx, def.CodeErr, err.Error(), err)
		return
	}
	helper.OKRsp(ctx, gin.H{})

}

func (ac *ArticleController) Delete(ctx *gin.Context) {
	params := &struct {
		Id int `form:"id" binding:"required"`
	}{}

	if err := ctx.ShouldBind(params); err != nil {
		helper.ErrRsp(ctx, def.CodeErr, "参数不正确", err)
		return
	}
	uid := ctx.GetInt("uid")
	as := service.NewArticleService(uid)

	err := as.Delete(params.Id)
	if err != nil {
		helper.ErrRsp(ctx, def.CodeErr, err.Error(), err)
		return
	}
	helper.OKRsp(ctx, gin.H{})
}

func (ac *ArticleController) GetDetail(ctx *gin.Context) {
	params := &struct {
		Id int `form:"id" binding:"required"`
	}{}

	if err := ctx.ShouldBind(params); err != nil {
		helper.ErrRsp(ctx, def.CodeErr, "参数不正确", err)
		return
	}
	uid := ctx.GetInt("uid")
	as := service.NewArticleService(uid)

	article, err := as.GetDetail(params.Id)
	if err != nil {
		helper.ErrRsp(ctx, def.CodeErr, err.Error(), err)
		return
	}
	helper.OKRsp(ctx, gin.H{
		"article": article,
	})
}

func (ac *ArticleController) List(ctx *gin.Context) {
	params := &struct {
		LastId int `form:"last_id" binding:"required"`
		Limit  int `form:"limit" binding:"required"`
	}{}

	if err := ctx.ShouldBind(params); err != nil {
		helper.ErrRsp(ctx, def.CodeErr, "参数不正确", err)
		return
	}
	uid := ctx.GetInt("uid")
	as := service.NewArticleService(uid)
	list, err := as.List(params.LastId, params.Limit)
	if err != nil {
		helper.ErrRsp(ctx, def.CodeErr, err.Error(), err)
		return
	}
	helper.OKRsp(ctx, gin.H{"list": list})
}

func (ac *ArticleController) ListMyArticle(ctx *gin.Context) {
	params := &struct {
		LastId int `form:"last_id" binding:"required"`
		Limit  int `form:"limit" binding:"required"`
	}{}

	if err := ctx.ShouldBind(params); err != nil {
		helper.ErrRsp(ctx, def.CodeErr, "参数不正确", err)
		return
	}
	uid := ctx.GetInt("uid")
	as := service.NewArticleService(uid)
	list, err := as.ListMyArticle(params.LastId, params.Limit)
	if err != nil {
		helper.ErrRsp(ctx, def.CodeErr, err.Error(), err)
		return
	}
	helper.OKRsp(ctx, gin.H{"list": list})
}
