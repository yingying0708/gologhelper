package gologhelper

import (
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

func PrintLogConsole(appname string, msg interface{}, log *logrus.Logger) {
	if msg != "" {
		pc, file, line, _ := runtime.Caller(2)
		f := runtime.FuncForPC(pc)
		hostname, err := os.Hostname()
		if err != nil {
			log.Println("获取hostname失败")
		}
		log.WithFields(logrus.Fields{
			"file":     file,
			"lineno":   line,
			"app_name": appname,
			"module":   strings.Split(f.Name(), ".")[0],
			"funcName": strings.Split(f.Name(), ".")[1],
			"log_time": time.Now().Format("2006-01-02 15:04:05"),
			"HOSTNAME": hostname,
		}).Println(msg)
	}
}
