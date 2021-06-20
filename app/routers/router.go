package routers

import (
	"yanglu/app/controller"
	"yanglu/interceptor"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	user(r)
	host(r)
	task(r)
	return r
}

func user(r *gin.Engine) {
	user := r.Group("/user")
	{
		user.POST("/add_user", interceptor.NewInterceptor().ParseToken, controller.NewUser().AddUser)
		user.POST("/login", controller.NewUser().Login)
		user.GET("/get_user_info", interceptor.NewInterceptor().ParseToken, controller.NewUser().GetUserInfo)
		user.POST("/set_authority", interceptor.NewInterceptor().ParseToken, controller.NewUser().SetAuthority)
		user.POST("/delete_user", interceptor.NewInterceptor().ParseToken, controller.NewUser().DeleteUser)
	}
}

func host(r *gin.Engine) {
	host := r.Group("host")
	{
		host.POST("/add", interceptor.NewInterceptor().ParseToken, controller.NewHost().Add)
		host.POST("/batch_add", interceptor.NewInterceptor().ParseToken, controller.NewHost().BatchAdd)
		host.POST("/update_department", interceptor.NewInterceptor().ParseToken, controller.NewHost().UpdateDepartment)
		host.POST("/search_host", interceptor.NewInterceptor().ParseToken, controller.NewHost().SearchHost)
		host.GET("/get_vulnerability_info", interceptor.NewInterceptor().ParseToken, controller.NewHost().GetVulnerabilityInfo)
	}
}

func task(r *gin.Engine) {
	task := r.Group("/task")
	{
		task.POST("/start_fast_task", interceptor.NewInterceptor().ParseToken, controller.NewTask().StartFastTask)
		task.GET("/get_progress", interceptor.NewInterceptor().ParseToken, controller.NewTask().GetProgress)
		task.POST("/add_timed_task", interceptor.NewInterceptor().ParseToken, controller.NewTask().AddTimedTask)
		task.POST("/add_repeat_task", interceptor.NewInterceptor().ParseToken, controller.NewTask().AddRepeatTask)
		task.GET("/get_detail", interceptor.NewInterceptor().ParseToken, controller.NewTask().GetDetail)
	}
}
