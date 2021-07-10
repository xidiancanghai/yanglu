package service

import (
	"yanglu/service/model"

	"github.com/sirupsen/logrus"
)

type ArticleService struct {
	uid     int
	article *model.ArticleInfo
}

func NewArticleService(uid int) *ArticleService {
	return &ArticleService{
		article: model.NewArticleInfo(),
		uid:     uid,
	}
}

func (as *ArticleService) Add(title string, tag string, content string, photos []string) error {
	as.article.Uid = as.uid
	as.article.Content.Title = title
	as.article.Content.Tag = tag
	as.article.Content.Content = content
	as.article.Content.Photos = photos
	err := as.article.Create()
	if err != nil {
		logrus.Error("Add err = ", err)
	}
	return err
}

func (as *ArticleService) List(lastId int, limit int) ([]*model.ArticleInfo, error) {
	return as.article.List(lastId, limit)
}

func (as *ArticleService) ListMyArticle(lastId int, limit int) ([]*model.ArticleInfo, error) {
	return as.article.ListMyArticle(as.uid, lastId, limit)
}
