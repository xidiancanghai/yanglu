package helper

import (
	"bytes"
	"errors"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func GetRandomStr(n int) string {
	letterRunes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var res bytes.Buffer
	rand.Seed(time.Now().Unix())
	for i := 0; i < n; i++ {
		index := rand.Intn(len(letterRunes))
		res.WriteByte(letterRunes[index])
	}
	return res.String()
}

func VerifyEmailFormat(email string) bool {
	//pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*` //匹配电子邮箱
	pattern := `^[0-9a-z][_.0-9a-z-]{0,31}@([0-9a-z][0-9a-z-]{0,30}[0-9a-z]\.){1,4}[a-z]{2,4}$`

	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

//mobile verify
func VerifyMobileFormat(mobileNum string) bool {
	regular := "^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199|(147))\\d{8}$"

	reg := regexp.MustCompile(regular)
	return reg.MatchString(mobileNum)
}

func StrToInts(str string) ([]int, error) {
	ss := strings.Split(str, ",")
	if len(ss) == 0 {
		return []int{}, errors.New("")
	}
	res := make([]int, len(ss))
	var err error = nil
	for k, v := range ss {
		res[k], err = strconv.Atoi(v)
		if err != nil {
			return res, err
		}
	}
	return res, nil
}
