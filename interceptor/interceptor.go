package interceptor

import (
	"net/http"
	"strconv"
	"time"
	"yanglu/config"
	"yanglu/def"
	"yanglu/helper"
	"yanglu/service"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Interceptor struct {
}

func NewInterceptor() *Interceptor {
	return &Interceptor{}
}

func (ic *Interceptor) LicenseExpired(ctx *gin.Context) {
	if config.LicenseInfoConf.ExpireTime < time.Now().Unix() {
		ctx.AbortWithStatusJSON(http.StatusOK, helper.Rsp(def.CodeErr, "当前license已经过期", nil))
		return
	}
}

func (ic *Interceptor) Cloud(ctx *gin.Context) {
	if !config.IsCloud() {
		ctx.AbortWithStatusJSON(http.StatusOK, helper.Rsp(def.CodeErr, "当前非云端版本", nil))
		return
	}
}

func (ic *Interceptor) ParseToken(ctx *gin.Context) {
	secret := def.ApiJwtSecretDev
	if config.IsOnline() {
		secret = def.ApiJwtSecret
	}
	token := ctx.GetHeader("token")
	if config.IsLocal() && token == "" {

		var uid int64 = 0
		if ctx.Request.Method == "POST" {
			uids := ctx.PostForm("uid")
			uid, _ = strconv.ParseInt(uids, 10, 64)
		} else {
			uid, _ = strconv.ParseInt(ctx.Query("uid"), 10, 64)
		}
		if uid != 0 {
			ctx.Set("uid", uid)
		}
		return
	}
	if !config.IsLocal() && len(token) <= 0 {
		ctx.AbortWithStatusJSON(http.StatusOK, helper.Rsp(def.CodeErr, "token无效", nil))
		return
	}
	claims, err := service.NewEmptyTokenService().CheckToken(token, secret)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"claims": claims,
			"err":    err,
		}).Warnln("ParseToken")
		ctx.AbortWithStatusJSON(http.StatusOK, helper.Rsp(def.CodeErr, "token无效", nil))
		return
	}
	ctx.Set("uid", claims.Uid)
}
