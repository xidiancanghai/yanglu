package service

import (
	"bytes"
	"errors"
	"io"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
	"yanglu/config"
	"yanglu/helper"

	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type UtilService struct {
}

func NewUtilService() *UtilService {
	return &UtilService{}
}

func (us *UtilService) GetCaptchaId() string {
	return captcha.NewLen(6)
}

func (us *UtilService) GetCaptcha(id string) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	if err := captcha.WriteImage(buf, id, 120, 80); err != nil {
		logrus.Error("GetCaptcha err ", err)
		return nil, err
	}
	return buf, nil
}

func (us *UtilService) ValidateCaptcha(id string, value string) error {
	if !config.IsCloud() {
		return nil
	}
	if config.IsLocal() && id == "" {
		return nil
	}
	if !config.IsCloud() {
		return nil
	}
	if captcha.VerifyString(id, value) {
		return nil
	}
	return errors.New("验证码错误或过期")
}

func (us *UtilService) UploadImage(file io.Reader, originFileName string) (string, error) {

	reg := regexp.MustCompile(`.+\.(jpg|bmp|gif|ico|pcx|jpeg|tif|png|raw|tga)$`)

	if !reg.MatchString(originFileName) {
		return "", errors.New("请输入正确的图片名称")
	}
	fileName, err := us.getFileName(originFileName)
	if err != nil {
		return "", err
	}
	ow, err := os.Create(fileName)
	if err != nil {
		return "", err
	}
	defer ow.Close()
	_, err = io.Copy(ow, file)
	if err != nil {
		logrus.Error("UploadImage err ", err)
		return "", err
	}
	ss := strings.Split(fileName, "/")
	return ss[len(ss)-1], nil
}

func (us *UtilService) getFileName(originFileName string) (string, error) {
	_, filename, _, _ := runtime.Caller(0)
	ss := strings.Split(filename, "yanglu")
	absolutelyPath := ss[0] + "yanglu/images"

	if err := us.CheckDir(absolutelyPath); err != nil {
		logrus.Error("getFileName err ", err)
		return "", err
	}
	ss = strings.Split(originFileName, ".")
	format := ss[len(ss)-1]
	return absolutelyPath + "/" + helper.GetRandomStr(10) + strconv.FormatInt(time.Now().UnixNano()/1000000, 10) + "." + format, nil

}

func (us *UtilService) CheckDir(absolutelyPath string) error {

	_, err := os.Stat(absolutelyPath) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return nil
		}

		return os.MkdirAll(absolutelyPath, os.ModePerm)

	}
	return nil
}

func (us *UtilService) DownloadImage(fileName string, ctx *gin.Context) error {
	_, filename, _, _ := runtime.Caller(0)
	ss := strings.Split(filename, "yanglu")
	absolutelyPath := ss[0] + "yanglu/images"
	fileName = absolutelyPath + "/" + fileName

	_, err := os.Stat(fileName)
	if err != nil && os.IsNotExist(err) {
		return errors.New("文件不存在")
	}
	ctx.Writer.Header().Add("Content-Type", "image/png")

	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(ctx.Writer, file)
	if err != nil {
		return err
	}
	return nil
}
