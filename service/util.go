package service

import (
	"bytes"
	"errors"
	"yanglu/config"

	"github.com/dchest/captcha"
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
	if captcha.VerifyString(id, value) {
		return nil
	}
	return errors.New("验证码错误或过期")
}
