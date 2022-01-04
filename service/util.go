package service

import (
	"bytes"
	"encoding/json"
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

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
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

func (s *UtilService) BuildPdf(list []map[string]interface{}) ([]byte, error) {

	if len(list) == 0 {
		return nil, errors.New("参数错误")
	}

	str := `
	<html>
		<head>
		<meta charset="utf-8">
		<title>检测结果</title>
		</head>
		<body>
		<table border="1">
	`
	keys := []string{}

	str += "<tr>"
	for k := range list[0] {
		keys = append(keys, k)
		str += "<td>" + k + "</td>"
	}
	str += "</tr>"

	for _, m := range list {

		str += "<tr>"
		for _, key := range keys {
			str += "<td>" + Strval(m[key]) + "</td>"
		}
		str += "</tr>"
	}
	str += `</table>
		</body>
		</html>
	`

	pdfg, err := wkhtmltopdf.NewPDFGenerator()

	if err != nil {
		return nil, err
	}

	pdfg.AddPage(wkhtmltopdf.NewPageReader(strings.NewReader(str)))

	if err = pdfg.Create(); err != nil {
		logrus.Error("Create err = ", err)
		return nil, err
	}

	return pdfg.Bytes(), nil
}

func Strval(value interface{}) string {
	// interface 转 string
	var key string
	if value == nil {
		return key
	}
	switch value.(type) {
	case float64:
		ft := value.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := value.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := value.(int)
		key = strconv.Itoa(it)
	case uint:
		it := value.(uint)
		key = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		key = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		key = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		key = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		key = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		key = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		key = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		key = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		key = strconv.FormatUint(it, 10)
	case string:
		key = value.(string)
	case []byte:
		key = string(value.([]byte))
	default:
		newValue, _ := json.Marshal(value)
		key = string(newValue)
	}

	return key
}
