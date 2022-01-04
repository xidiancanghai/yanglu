package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
	"yanglu/app/routers"
	"yanglu/config"
	"yanglu/service"
	"yanglu/service/data"
	"yanglu/service/logger"

	"github.com/facebookgo/grace/gracehttp"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron"
)

func main() {

	// read toml file to init config, must run this first.

	// init redis
	//data.NewPackRedis()

	config.InitEnvConf()
	config.InitGoodsConif()
	data.InitMysql()
	data.InitMemoryCache()

	logger.InitLogger(config.GetLogPath(), nil)
	if err := config.InitLicenseConfig(); err != nil {
		log.Fatal("启动错误 err = ", err)
	}

	port := strconv.Itoa(config.GetHttpPort())
	flag.Parse()

	//路由配置
	r := routers.InitRouter()

	//r.HTMLRender, _ = pongo2gin.NewRenderer(adminRoot + "dist")

	log.Println("服务正在启动，监听端口:", port, ",PID:", strconv.Itoa(os.Getpid()), gin.Version, "version: temp")
	server := &http.Server{
		Addr:         ":" + port,
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 20 * time.Second,
		Handler:      r,
	}
	c := cron.New()
	c.AddFunc("@every 30s", service.NewSmartTask().TaskCheck)
	c.Start()

	err := gracehttp.Serve(server)
	if err != nil {
		log.Fatal("服务启动失败:", err.Error())
	}
}
