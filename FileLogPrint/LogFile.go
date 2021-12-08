package loghelper

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
type LogFile struct {
	File     string      `json:"file"`
	LineNo   int         `json:"lineno"`
	App_Name string      `json:"app_name"`
	Module   string      `json:"module"`
	FuncName string      `json:"funcName"`
	Log_Time string      `json:"log_time"`
	HOSTNAME string      `json:"hostname"`
	Msg      interface{} `json:"msg"`
}

func PrintLogFile(writer *rotatelogs.RotateLogs, appname string, msg interface{}, log *logrus.Logger) {
	pc, file, line, _ := runtime.Caller(2)
	f := runtime.FuncForPC(pc)
	hostname, err := os.Hostname()
	if err != nil {
		log.Println("获取hostname失败")
	}
	module := strings.Split(f.Name(), ".")[0]
	funcName := strings.Split(f.Name(), ".")[1]
	log_time := time.Now().Format("2006-01-02 15:04:05")

	entity := LogFile{
		File:     file,
		FuncName: funcName,
		LineNo:   line,
		App_Name: appname,
		Module:   module,
		Log_Time: log_time,
		HOSTNAME: hostname,
		Msg:      msg,
	}

	if res, err := json.Marshal(&entity); err == nil {
		writer.Write(res)
		writer.Write([]byte("\n"))
	}
	log.SetOutput(writer)
}
