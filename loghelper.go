package gologhelper

import (
	"encoding/json"
	"strings"
	"time"

	consolehelper "github.com/yingying0708/gologhelper/ConsoleLogPrint"
	filehelper "github.com/yingying0708/gologhelper/FileLogPrint"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

var logs = logrus.New()

//初始化
func init() {
	//设置日志格式为json
	logs.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:03:04",
	})

}

// 日志帮助类
type LogHelper struct {
	AppName      string
	ConsolePrint bool
	LogPath      string
	BackupCount  int
	When         string
	LogLevel     string
	TraceWriter  *rotatelogs.RotateLogs
	DebugWriter  *rotatelogs.RotateLogs
	InfoWriter   *rotatelogs.RotateLogs
	WarnWriter   *rotatelogs.RotateLogs
	ErrorWriter  *rotatelogs.RotateLogs
}

// 日志级别
const (
	Log_Trace = iota
	Log_Debug
	Log_Info
	Log_Warn
	Log_Error
)

// 根据输入的日志级别，返回匹配的自定义常数
func GetLogLevel(loglevel string) int {
	res := Log_Error
	if loglevel == "info" {
		res = Log_Info
	}
	if loglevel == "warn" {
		res = Log_Warn
	}
	if loglevel == "debug" {
		res = Log_Debug
	}
	if loglevel == "trace" {
		res = Log_Trace
	}
	if loglevel == "error" {
		res = Log_Error
	}
	return res
}

//外部调用函数生成日志类
func GetLogHelper(app_name, log_path string) *LogHelper {
	logPath := log_path + app_name + "_p1_error.log"
	return &LogHelper{
		AppName:      app_name,
		LogPath:      log_path,
		ConsolePrint: false,
		BackupCount:  15,
		When:         "D",
		LogLevel:     "error",
		ErrorWriter:  getWriter(logPath, "D", 15),
	}
}

//设置writer
func (log *LogHelper) SetWriter() {
	if GetLogLevel(log.LogLevel) <= Log_Trace {
		logPath := GetLogPath(log, "trace")
		log.TraceWriter = getWriter(logPath, log.When, log.BackupCount)
	}
	if GetLogLevel(log.LogLevel) <= Log_Debug {
		logPath := GetLogPath(log, "debug")
		log.DebugWriter = getWriter(logPath, log.When, log.BackupCount)
	}
	if GetLogLevel(log.LogLevel) <= Log_Info {
		logPath := GetLogPath(log, "info")
		log.InfoWriter = getWriter(logPath, log.When, log.BackupCount)
	}
	if GetLogLevel(log.LogLevel) <= Log_Warn {
		logPath := GetLogPath(log, "warn")
		log.WarnWriter = getWriter(logPath, log.When, log.BackupCount)
	}
	if GetLogLevel(log.LogLevel) <= Log_Error {
		logPath := GetLogPath(log, "error")
		log.ErrorWriter = getWriter(logPath, log.When, log.BackupCount)
	}
}

//返回路径
func GetLogPath(log *LogHelper, levelStr string) string {
	logPath := log.LogPath + log.AppName + "_p1_" + levelStr + ".log"
	return logPath
}

//设置when(D:天，H：小时，M：分钟，默认是D)
func (log *LogHelper) SetWhen(when string) *LogHelper {
	if when != "" {
		log.When = strings.ToUpper(when)
		log.SetWriter()
	}
	return log
}

//设置日志级别（error,debug,info,trace,warn 默认是error）
func (log *LogHelper) SetLogLevel(level string) *LogHelper {
	levelstr := "error"
	if level != "" {
		levelstr = strings.ToLower(level)
	}
	log.LogLevel = levelstr
	log.SetWriter()
	return log
}

//设置是否控制台打印默认是false
func (log *LogHelper) SetConsolePrint(isPrint bool) *LogHelper {
	log.ConsolePrint = isPrint
	return log
}

//设置多少个文件后进行回滚操作默认是15
func (log *LogHelper) SetBackupCount(backupCount int) *LogHelper {
	if backupCount > 0 {
		log.BackupCount = backupCount
		log.SetWriter()
	}
	return log
}

//Trace
func (log *LogHelper) Trace(param ...interface{}) {
	msg, extra := getParams(param...)
	if msg != nil && extra != nil {
		log.traceCustom(msg, extra)
	} else {
		log.trace(msg)
	}
}

//Debug
func (log *LogHelper) Debug(param ...interface{}) {
	msg, extra := getParams(param...)
	if msg != nil && extra != nil {
		log.debugCustom(msg, extra)
	} else {
		log.debug(msg)
	}
}

//Info
func (log *LogHelper) Info(param ...interface{}) {
	msg, extra := getParams(param...)
	if msg != nil && extra != nil {
		log.infoCustom(msg, extra)
	} else {
		log.info(msg)
	}
}

//Warn
func (log *LogHelper) Warn(param ...interface{}) {
	msg, extra := getParams(param...)
	if msg != nil && extra != nil {
		log.warnCustom(msg, extra)
	} else {
		log.warn(msg)
	}
}

//Error
func (log *LogHelper) Error(param ...interface{}) {
	msg, extra := getParams(param...)
	if msg != nil && extra != nil {
		log.errorCustom(msg, extra)
	} else {
		log.error(msg)
	}
}

//拆分参数
func getParams(param ...interface{}) (interface{}, map[string]interface{}) {
	var msg interface{}
	if len(param) == 1 {
		msg = param[0]
		return msg, nil
	}
	if len(param) == 2 {
		var mapResult map[string]interface{}
		msg = param[0]
		p1 := param[1]
		if extra, ok := p1.(string); ok {
			//判断extra的前几位是不是extra={}
			prefix := extra[0:6]
			if prefix == "extra=" {
				jsonContent := extra[6:]
				if err := json.Unmarshal([]byte(jsonContent), &mapResult); err != nil {
					mapResult = make(map[string]interface{}, 1)
					mapResult["error_msg"] = "extra数据格式错误"
				}
			}
		}
		return msg, mapResult
	}
	return nil, nil
}

//实例化rotatelogs.RotateLogs
func getWriter(logPath, When string, backupCount int) *rotatelogs.RotateLogs {
	//路径
	writer, _ := rotatelogs.New(
		logPath+".%Y%m%d.log",
		rotatelogs.WithRotationCount(uint(backupCount)),
		rotatelogs.WithRotationTime(time.Duration(24)*time.Hour),
	)
	if When == "H" {
		writer, _ = rotatelogs.New(
			logPath+".%Y%m%d%H.log",
			rotatelogs.WithRotationCount(uint(backupCount)),
			rotatelogs.WithRotationTime(time.Duration(60)*time.Minute),
		)
	}
	if When == "M" {
		writer, _ = rotatelogs.New(
			logPath+".%Y%m%d%H%M.log",
			rotatelogs.WithRotationCount(uint(backupCount)),
			rotatelogs.WithRotationTime(time.Duration(60)*time.Second),
		)
	}
	return writer
}

func (log *LogHelper) info(msg interface{}) {
	if GetLogLevel(log.LogLevel) <= Log_Info {
		if log.ConsolePrint {
			consolehelper.PrintLogConsole(log.AppName, "INFO", msg, logs)
		}
		filehelper.PrintLogFile(log.InfoWriter, log.AppName, "INFO", msg, logs)
	}
}

func (log *LogHelper) trace(msg interface{}) {
	if GetLogLevel(log.LogLevel) <= Log_Trace {
		if log.ConsolePrint {
			consolehelper.PrintLogConsole(log.AppName, "TRACE", msg, logs)
		}
		filehelper.PrintLogFile(log.TraceWriter, log.AppName, "TRACE", msg, logs)
	}
}

func (log *LogHelper) debug(msg interface{}) {
	if GetLogLevel(log.LogLevel) <= Log_Debug {
		if log.ConsolePrint {
			consolehelper.PrintLogConsole(log.AppName, "DEBUG", msg, logs)
		}
		filehelper.PrintLogFile(log.DebugWriter, log.AppName, "DEBUG", msg, logs)
	}
}

func (log *LogHelper) warn(msg interface{}) {
	if GetLogLevel(log.LogLevel) <= Log_Warn {
		if log.ConsolePrint {
			consolehelper.PrintLogConsole(log.AppName, "WARN", msg, logs)
		}
		filehelper.PrintLogFile(log.WarnWriter, log.AppName, "WARN", msg, logs)
	}
}

func (log *LogHelper) error(msg interface{}) {
	if GetLogLevel(log.LogLevel) <= Log_Error {
		if log.ConsolePrint {
			consolehelper.PrintLogConsole(log.AppName, "ERROR", msg, logs)
		}
		filehelper.PrintLogFile(log.ErrorWriter, log.AppName, "ERROR", msg, logs)
	}
}

func (log *LogHelper) infoCustom(msg interface{}, fields map[string]interface{}) {
	if GetLogLevel(log.LogLevel) <= Log_Info {
		if log.ConsolePrint {
			consolehelper.PrintLogConsoleCustom(log.AppName, "INFO", msg, fields, logs)
		}
		filehelper.PrintLogFileCustom(log.InfoWriter, log.AppName, "INFO", msg, fields, logs)
	}
}

func (log *LogHelper) traceCustom(msg interface{}, fields map[string]interface{}) {
	if GetLogLevel(log.LogLevel) <= Log_Trace {
		if log.ConsolePrint {
			consolehelper.PrintLogConsoleCustom(log.AppName, "TRACE", msg, fields, logs)
		}
		filehelper.PrintLogFileCustom(log.TraceWriter, log.AppName, "TRACE", msg, fields, logs)
	}
}

func (log *LogHelper) debugCustom(msg interface{}, fields map[string]interface{}) {
	if GetLogLevel(log.LogLevel) <= Log_Debug {
		if log.ConsolePrint {
			consolehelper.PrintLogConsoleCustom(log.AppName, "DEBUG", msg, fields, logs)
		}
		filehelper.PrintLogFileCustom(log.DebugWriter, log.AppName, "DEBUG", msg, fields, logs)
	}
}

func (log *LogHelper) warnCustom(msg interface{}, fields map[string]interface{}) {
	if GetLogLevel(log.LogLevel) <= Log_Warn {
		if log.ConsolePrint {
			consolehelper.PrintLogConsoleCustom(log.AppName, "WARN", msg, fields, logs)
		}
		filehelper.PrintLogFileCustom(log.WarnWriter, log.AppName, "WARN", msg, fields, logs)
	}
}

func (log *LogHelper) errorCustom(msg interface{}, fields map[string]interface{}) {
	if GetLogLevel(log.LogLevel) <= Log_Error {
		if log.ConsolePrint {
			consolehelper.PrintLogConsoleCustom(log.AppName, "ERROR", msg, fields, logs)
		}
		filehelper.PrintLogFileCustom(log.ErrorWriter, log.AppName, "ERROR", msg, fields, logs)
	}
}
