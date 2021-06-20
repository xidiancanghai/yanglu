package interceptor

import (
	"net/http"
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

func (ic *Interceptor) ParseToken(ctx *gin.Context) {
	secret := def.ApiJwtSecretDev
	if config.IsOnline() {
		secret = def.ApiJwtSecret
	}
	token := ctx.GetHeader("token")
	if config.IsLocal() {
		ctx.Set("uid", 1)
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
