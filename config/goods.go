package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path"
	"runtime"
)

type GoodsConfig struct {
	OndeNodePrice int   `json:"one_node_price"`
	OneMonthPrice int   `json:"one_month_price"`
	NodeConf      []int `json:"node_conf"`
	TimeConf      []int `json:"time_conf"`
}

var goodsConfigObj *GoodsConfig

func InitGoodsConif() {

	_, filename, _, _ := runtime.Caller(0)

	filePath := path.Dir(filename) + "/goods.json"

	data, err := ioutil.ReadFile(filePath)

	if err != nil {
		log.Fatal("找不到价格配置文件")
		return
	}

	goodsConfigObj = new(GoodsConfig)

	err = json.Unmarshal(data, goodsConfigObj)

	if err != nil {
		log.Fatal("配置文件解析错误")
		return
	}
}

func GetGoodsConfig() *GoodsConfig {
	return goodsConfigObj
}
