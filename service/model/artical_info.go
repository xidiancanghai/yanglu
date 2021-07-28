package model

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strconv"
	"time"
	"yanglu/service/data"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ArticleContent struct {
	Title   string   `json:"title"`
	Tag     string   `json:"tag"`
	Content string   `json:"content"`
	Photos  []string `json:"photos"`
}

func (a ArticleContent) Value() (driver.Value, error) {
	b, err := json.Marshal(a)
	return string(b), err
}

func (a *ArticleContent) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), a)
}

type ArticleInfo struct {
	Id         int            `json:"id"`
	Uid        int            `json:"uid"`
	Content    ArticleContent `json:"content"`
	IsDelete   int            `json:"-"`
	UpdateTime int64          `json:"-"`
	CreateTime int64          `json:"create_time"`
}

func NewArticleInfo() *ArticleInfo {
	return &ArticleInfo{}
}

func (a *ArticleInfo) TableName() string {
	return "article_info"
}

func (a *ArticleInfo) Create() error {
	if a.Uid == 0 {
		return errors.New("参数错误")
	}
	if a.CreateTime == 0 {
		a.CreateTime = time.Now().Unix()
	}
	a.UpdateTime = a.CreateTime
	tx := data.GetDB().Create(a)
	if tx.Error != nil {
		logrus.Error("Create err ", tx)
	}
	return tx.Error
}

func (a *ArticleInfo) GetArticle(id int) (*ArticleInfo, error) {

	if id == 0 {
		return nil, errors.New("id错误")
	}
	res := new(ArticleInfo)
	tx := data.GetDB().Where(" id = ? and is_delete = 0", id).First(res)
	if tx.Error != nil && tx.Error != gorm.ErrRecordNotFound {
		logrus.Error("GetArticle err tx = ", tx.Error)
		return nil, tx.Error
	}
	return res, nil
}

func (a *ArticleInfo) Updates(m map[string]interface{}) error {
	if len(m) == 0 {
		return errors.New("参数错误")
	}
	if a.Id == 0 {
		return errors.New("主键id错误")
	}
	tx := data.GetDB().Model(a).Updates(m)
	if tx.Error != nil {
		logrus.Error("Updates err ", tx)
		return tx.Error
	}
	return nil
}

func (a *ArticleInfo) List(lastId int, limit int) ([]*ArticleInfo, error) {
	sqll := ""
	if lastId == -1 {
		sqll = "select id, uid, content, create_time from " + a.TableName() + " where is_delete = 0 order by id desc limit ?"
	} else {
		sqll = "select id, uid, content, create_time from " + a.TableName() + " where id < " + strconv.Itoa(lastId) + " and is_delete = 0 order by id desc limit ?"
	}
	list := []*ArticleInfo{}
	rows, err := data.GetDB().Raw(sqll, limit).Rows()
	if err != nil && err != sql.ErrNoRows {
		logrus.Error("List err = ", err)
		return nil, err
	}
	if err != nil {
		return list, nil
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var uid int
		var content string
		var createTime int64
		rows.Scan(&id, &uid, &content, &createTime)
		ac := new(ArticleContent)
		json.Unmarshal([]byte(content), ac)
		if ac.Photos == nil {
			ac.Photos = []string{}
		}
		list = append(list, &ArticleInfo{
			Id:         id,
			Uid:        uid,
			Content:    *ac,
			CreateTime: createTime,
		})
	}
	return list, nil
}

func (a *ArticleInfo) ListMyArticle(uid int, lastId int, limit int) ([]*ArticleInfo, error) {
	sqll := ""
	if lastId == -1 {
		sqll = "select id, uid, content, create_time from " + a.TableName() + " where uid = ? and is_delete = 0 order by id desc limit ?"
	} else {
		sqll = "select id, uid, content, create_time from " + a.TableName() + " where uid = ? and id < " + strconv.Itoa(lastId) + " and is_delete = 0 order by id desc limit ?"
	}
	list := []*ArticleInfo{}
	rows, err := data.GetDB().Raw(sqll, uid, limit).Rows()
	if err != nil && err != sql.ErrNoRows {
		logrus.Error("List err = ", err)
		return nil, err
	}
	if err != nil {
		return list, nil
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var uid int
		var content string
		var createTime int64
		rows.Scan(&id, &uid, &content, &createTime)
		ac := new(ArticleContent)
		json.Unmarshal([]byte(content), ac)
		if ac.Photos == nil {
			ac.Photos = []string{}
		}
		list = append(list, &ArticleInfo{
			Id:         id,
			Uid:        uid,
			Content:    *ac,
			CreateTime: createTime,
		})
	}
	return list, nil
}
