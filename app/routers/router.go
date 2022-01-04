package routers

import (
	"net/http"
	"yanglu/app/controller"
	"yanglu/interceptor"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.New()

	adminRoot := "../"
	r.StaticFS("/css", http.Dir(adminRoot+"dist/css"))
	r.StaticFS("/js", http.Dir(adminRoot+"dist/js"))
	r.StaticFS("/img", http.Dir(adminRoot+"dist/img"))
	r.StaticFS("/fonts", http.Dir(adminRoot+"dist/fonts"))
	r.StaticFile("/favicon.ico", adminRoot+"dist/favicon.ico")
	r.HTMLRender = New(adminRoot)

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "/dist/index.html", pongo2.Context{})
	})

	r.Use(gin.Recovery())
	r.Use(interceptor.NewInterceptor().LicenseExpired)
	r.Use()
	user(r)
	host(r)
	task(r)
	actionLog(r)
	config(r)
	util(r)
	order(r)
	article(r)
	return r
}

func user(r *gin.Engine) {
	user := r.Group("/user")
	{
		user.POST("/register", interceptor.NewInterceptor().Cloud, controller.NewUser().Register)
		user.POST("/add_user", interceptor.NewInterceptor().ParseToken, controller.NewUser().AddUser)
		user.POST("/reset_passwd", interceptor.NewInterceptor().ParseToken, controller.NewUser().ResetPasswd)
		user.POST("/login", controller.NewUser().Login)
		user.POST("/logout", controller.NewUser().Logout)
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
		host.GET("/get_vulnerability_pdf", interceptor.NewInterceptor().ParseToken, controller.NewHost().GetVulnerabilityPdf)
		host.GET("/list_all", interceptor.NewInterceptor().ParseToken, controller.NewHost().ListAll)
		host.GET("/vulnerability_distribute", interceptor.NewInterceptor().ParseToken, controller.NewHost().VulnerabilityDistribute)
		host.GET("/system_os_distribute", interceptor.NewInterceptor().ParseToken, controller.NewHost().SystemOsDistribute)
		host.POST("/set_ip", interceptor.NewInterceptor().ParseToken, controller.NewHost().SetIp)
		host.POST("/delete", interceptor.NewInterceptor().ParseToken, controller.NewHost().Delete)
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
		log.GET("/search_log", interceptor.NewInterceptor().ParseToken, controller.NewLog().SearchLog)
	}

}

func util(r *gin.Engine) {
	util := r.Group("/util")
	{
		util.GET("/get_captcha", controller.NewUtilController().GetCaptcha)
		util.GET("/get_captcha_id", controller.NewUtilController().GetCaptchaId)
		util.GET("/get_system_info", controller.NewUtilController().GetSystemInfo)
		util.POST("/upload_images", controller.NewUtilController().UploadImages)
		util.GET("/download_images", controller.NewUtilController().DownloadImage)
	}
}

func order(r *gin.Engine) {
	order := r.Group("/order")
	{
		order.GET("/config", controller.NewOrder().GetConfig)
		order.POST("/create", controller.NewOrder().Create)
	}
}

func article(r *gin.Engine) {
	article := r.Group("/article", interceptor.NewInterceptor().ParseToken)
	{
		article.POST("/add", controller.NewArticleController().Add)
		article.GET("/list", controller.NewArticleController().List)
		article.GET("/list_my_article", controller.NewArticleController().ListMyArticle)
		article.POST("/delete", controller.NewArticleController().Delete)
		article.GET("/get_detail", controller.NewArticleController().GetDetail)
	}
}
