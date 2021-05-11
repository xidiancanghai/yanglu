package routers

import (
	"yanglu/app/controller"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	user(r)
	return r
}

func user(r *gin.Engine) {
	user := r.Group("/user")
	{
		user.POST("/add_user", controller.NewUser().AddUser)
		user.POST("/login", controller.NewUser().login)
	}
}
