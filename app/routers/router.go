package routers

import (
	"yanglu/app/controller"
	"yanglu/interceptor"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(interceptor.NewInterceptor().LicenseExpired)
	r.Use()
	user(r)
	host(r)
	task(r)
	actionLog(r)
	config(r)
	util(r)
	return r
}

func user(r *gin.Engine) {
	user := r.Group("/user")
	{
		user.POST("/user/register", interceptor.NewInterceptor().Cloud, controller.NewUser().Register)
		user.POST("/add_user", interceptor.NewInterceptor().ParseToken, controller.NewUser().AddUser)
		user.POST("/reset_passwd", interceptor.NewInterceptor().ParseToken, controller.NewUser().ResetPasswd)
		user.POST("/login", controller.NewUser().Login)
		user.GET("/get_user_info", interceptor.NewInterceptor().ParseToken, controller.NewUser().GetUserInfo)
		user.POST("/set_authority", interceptor.NewInterceptor().ParseToken, controller.NewUser().SetAuthority)
		user.POST("/delete_user", interceptor.NewInterceptor().ParseToken, controller.NewUser().DeleteUser)
		user.GET("/list_users", interceptor.NewInterceptor().ParseToken, controller.NewUser().ListUsers)
		user.POST("/find_passwd", controller.NewUser().FindPassWd)
	}
}

func config(r *gin.Engine) {
	conf := r.Group("/config")
	{
		conf.GET("/get_const_config", controller.NewConfController().GetConf)
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
		host.GET("/list_all", interceptor.NewInterceptor().ParseToken, controller.NewHost().ListAll)
		host.GET("/vulnerability_distribute", interceptor.NewInterceptor().ParseToken, controller.NewHost().VulnerabilityDistribute)
		host.GET("/system_os_distribute", interceptor.NewInterceptor().ParseToken, controller.NewHost().SystemOsDistribute)
		host.POST("/set_ip", interceptor.NewInterceptor().ParseToken, controller.NewHost().SetIp)
	}
}

func task(r *gin.Engine) {
	task := r.Group("/task")
	{
		task.POST("/start_fast_task", interceptor.NewInterceptor().ParseToken, controller.NewTask().StartFastTask)
		task.GET("/get_progress", interceptor.NewInterceptor().ParseToken, controller.NewTask().GetProgress)
		task.POST("/add_timed_task", interceptor.NewInterceptor().ParseToken, controller.NewTask().AddTimedTask)
		task.POST("/add_repeat_task", interceptor.NewInterceptor().ParseToken, controller.NewTask().AddRepeatTask)
		task.GET("/curl_task_info", interceptor.NewInterceptor().ParseToken, controller.NewTask().GetDetail)
		task.GET("/check_info", interceptor.NewInterceptor().ParseToken, controller.NewTask().CheckInfo)
	}
}

func actionLog(r *gin.Engine) {
	log := r.Group("/log")
	{
		log.GET("/list", interceptor.NewInterceptor().ParseToken, controller.NewLog().List)
	}

}

func util(r *gin.Engine) {
	util := r.Group("/util")
	{
		util.GET("/get_captcha", controller.NewUtilController().GetCaptcha)
		util.GET("/get_captcha_id", controller.NewUtilController().GetCaptchaId)
		util.GET("/get_system_info", controller.NewUtilController().GetSystemInfo)
	}
}
