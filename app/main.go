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
	"yanglu/service/data"

	"github.com/facebookgo/grace/gracehttp"
	"github.com/gin-gonic/gin"
)

func main() {

	// read toml file to init config, must run this first.

	// init redis
	//data.NewPackRedis()
	config.InitEnvConf()
	data.InitMysql()

	port := flag.String("port", "9080", "port")
	flag.Parse()

	//路由配置
	r := routers.InitRouter()

	log.Println("服务正在启动，监听端口:", *port, ",PID:", strconv.Itoa(os.Getpid()), gin.Version, "version: temp")
	server := &http.Server{
		Addr:         ":" + *port,
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 20 * time.Second,
		Handler:      r,
	}
	err := gracehttp.Serve(server)
	if err != nil {
		log.Fatal("服务启动失败:", err.Error())
	}
}
