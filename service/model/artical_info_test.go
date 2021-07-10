package model

import (
	"fmt"
	"testing"
	"yanglu/config"
	"yanglu/service/data"
	"yanglu/service/logger"
)

func TestArticleInfo(t *testing.T) {

	config.InitEnvConf("../../config/env.toml")

	// init log
	logger.InitLogger("", nil)

	// init db
	data.InitMysql()

	// a := ArticleInfo{
	// 	Uid: 1,
	// 	Content: ArticleContent{
	// 		Title:   "测试",
	// 		Tag:     "你好啊",
	// 		Content: "测试啊",
	// 	},
	// }
	// a.Create()

	a := NewArticleInfo()
	list, err := a.List(2, 10)
	fmt.Println(err)
	if err != nil {
		return
	}
	for _, v := range list {
		fmt.Println("v = ", *v)
	}
}
