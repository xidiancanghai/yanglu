package logger

import (
	"io/ioutil"
	"log"
	"os"
	"time"
	"yanglu/config"

	rotate "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

func InitLogger(logPath string, prefix map[string]interface{}) {
	if len(logPath) == 0 {
		logPath = config.GetLogPath()
	}

	exist, _ := pathExists(logPath)
	if !exist {
		err := os.MkdirAll(logPath, os.ModePerm)
		if err != nil {
			log.Fatalf("failed to mkdir, path: %v, err: %v", logPath, err)
		}
	}
	logrus.SetOutput(ioutil.Discard) //使用hook输出日志，丢弃原有的write操作

	logrus.SetLevel(logrus.InfoLevel)
	if config.IsTest() {
		logrus.SetLevel(logrus.DebugLevel)
	}

	if config.IsLocal() {
		logrus.SetOutput(os.Stdout)
	}

	logMain, err := rotate.New(logPath+"/main.log.%Y%m%d%H", rotate.WithMaxAge(30*24*time.Hour), rotate.WithRotationTime(time.Hour))
	if err != nil {
		log.Fatalf("failed to rotate logs, err: %v", err)
	}

	logError, err := rotate.New(logPath+"/error.log.%Y%m%d%H", rotate.WithMaxAge(30*24*time.Hour), rotate.WithRotationTime(time.Hour))
	if err != nil {
		log.Fatalf("failed to rotate logs, err: %v", err)
	}
	//为不同级别设置不同的输出目的
	lfHook := NewHook(
		WriterMap{
			logrus.DebugLevel: logMain,
			logrus.InfoLevel:  logMain,
			logrus.WarnLevel:  logMain,
			logrus.ErrorLevel: logError,
			logrus.FatalLevel: logError,
			logrus.PanicLevel: logError,
		},
		&logrus.JSONFormatter{},
	)
	logrus.AddHook(lfHook)
	//logrus.SetReportCaller(true)

}

func LoggerTemp(logPath string, prefix map[string]interface{}) *logrus.Logger {
	if len(logPath) == 0 {
		logPath = config.GetLogPath()
	}

	exist, _ := pathExists(logPath)
	if !exist {
		err := os.MkdirAll(logPath, os.ModePerm)
		if err != nil {
			log.Fatalf("failed to mkdir, path: %v, err: %v", logPath, err)
		}
	}

	sLog := logrus.New()

	sLog.SetOutput(ioutil.Discard) //使用hook输出日志，丢弃原有的write操作

	sLog.SetLevel(logrus.InfoLevel)
	if config.IsTest() {
		sLog.SetLevel(logrus.DebugLevel)
	}

	if config.IsLocal() {
		sLog.SetOutput(os.Stdout)
	}

	logMain, err := rotate.New(logPath+"/main.log.%Y%m%d%H", rotate.WithMaxAge(30*24*time.Hour), rotate.WithRotationTime(time.Hour))
	if err != nil {
		log.Fatalf("failed to rotate logs, err: %v", err)
	}

	logError, err := rotate.New(logPath+"/error.log.%Y%m%d%H", rotate.WithMaxAge(30*24*time.Hour), rotate.WithRotationTime(time.Hour))
	if err != nil {
		log.Fatalf("failed to rotate logs, err: %v", err)
	}
	//为不同级别设置不同的输出目的

	// lfHook := lfshook.NewHook(
	// 	lfshook.WriterMap{
	// 		logrus.DebugLevel: logMain,
	// 		logrus.InfoLevel:  logMain,
	// 		logrus.WarnLevel:  logMain,
	// 		logrus.ErrorLevel: logError,
	// 		logrus.FatalLevel: logError,
	// 		logrus.PanicLevel: logError,
	// 	},
	// 	&logrus.JSONFormatter{},
	// )
	lfHook := NewHook(
		WriterMap{
			logrus.DebugLevel: logMain,
			logrus.InfoLevel:  logMain,
			logrus.WarnLevel:  logMain,
			logrus.ErrorLevel: logError,
			logrus.FatalLevel: logError,
			logrus.PanicLevel: logError,
		},
		&logrus.JSONFormatter{},
	)
	sLog.AddHook(lfHook)
	//logrus.SetReportCaller(true)
	//std = sLog.WithField("srv", "api")
	//
	////std = logrus.WithField("service", "prefix")
	//if prefix == nil || len(prefix) == 0 {
	//	std = logrus.WithField("prefix", "nil")
	//	prefix = make(map[string]interface{})
	//}
	//for k, val := range prefix {
	//	std = std.WithField(k, val)
	//}
	return sLog
}

// 判断文件夹是否存在
func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
