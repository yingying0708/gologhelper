package gologhelper

import (
	"encoding/json"
	"os"
	"runtime"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

//打印的数据结构体
type LogFileCustom struct {
	File     string      `json:"file"`
	LineNo   int         `json:"lineno"`
	App_Name string      `json:"app_name"`
	Module   string      `json:"module"`
	FuncName string      `json:"funcName"`
	Log_Time string      `json:"log_time"`
	HOSTNAME string      `json:"hostname"`
	Level    string      `json:"level"`
	Msg      interface{} `json:"msg"`
}

func PrintLogFileCustom(writer *rotatelogs.RotateLogs, appname, level string, fields map[string]interface{}, log *logrus.Logger) {
	pc, file, line, _ := runtime.Caller(2)
	f := runtime.FuncForPC(pc)
	hostname, err := os.Hostname()
	if err != nil {
		log.Println("获取hostname失败")
	}
	module := strings.Split(f.Name(), ".")[0]
	funcName := strings.Split(f.Name(), ".")[1]
	log_time := time.Now().Format("2006-01-02 15:04:05")

	fields["file"] = file
	fields["lineno"] = line
	fields["app_name"] = appname
	fields["module"] = module
	fields["funcName"] = funcName
	fields["log_time"] = log_time
	fields["hostname"] = hostname
	fields["level"] = level

	if res, err := json.Marshal(fields); err == nil {
		writer.Write(res)
		writer.Write([]byte("\n"))
	}
	log.SetOutput(writer)
}
