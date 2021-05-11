package helper

import (
	"net/http"
	"yanglu/def"

	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func Rsp(code int, msg string, data interface{}) *Response {
	if code != 0 {
		logrus.WithField("data", data).WithField("msg", msg).WithField("code", code).Debug("failed rsp")
	}
	rsp := &Response{
		Code: code, Message: msg, Data: data,
	}
	return rsp
}

func ErrRsp(ctx *gin.Context, code int, msg string, err error) {
	logrus.WithError(err).WithField("msg", msg).Warnf("failed request: %v", ctx.Request.RequestURI)
	if err != nil {
		ctx.JSON(http.StatusOK, Rsp(code, msg, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, Rsp(code, msg, nil))
}

func OKRsp(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, Rsp(def.CodeOK, "ok", data))
}
