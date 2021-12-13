package gologhelper

import (
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

func (log *LogHelper) Info(msg interface{}) {
	if GetLogLevel(log.LogLevel) <= Log_Info {
		if log.ConsolePrint {
			consolehelper.PrintLogConsole(log.AppName, "info", msg, logs)
		}
		filehelper.PrintLogFile(log.InfoWriter, log.AppName, "info", msg, logs)
	}
}

func (log *LogHelper) Trace(msg interface{}) {
	if GetLogLevel(log.LogLevel) <= Log_Trace {
		if log.ConsolePrint {
			consolehelper.PrintLogConsole(log.AppName, "trace", msg, logs)
		}
		filehelper.PrintLogFile(log.TraceWriter, log.AppName, "trace", msg, logs)
	}
}

func (log *LogHelper) Debug(msg interface{}) {
	if GetLogLevel(log.LogLevel) <= Log_Debug {
		if log.ConsolePrint {
			consolehelper.PrintLogConsole(log.AppName, "debug", msg, logs)
		}
		filehelper.PrintLogFile(log.DebugWriter, log.AppName, "debug", msg, logs)
	}
}

func (log *LogHelper) Warn(msg interface{}) {
	if GetLogLevel(log.LogLevel) <= Log_Warn {
		if log.ConsolePrint {
			consolehelper.PrintLogConsole(log.AppName, "warn", msg, logs)
		}
		filehelper.PrintLogFile(log.WarnWriter, log.AppName, "warn", msg, logs)
	}
}

func (log *LogHelper) Error(msg interface{}) {
	if GetLogLevel(log.LogLevel) <= Log_Error {
		if log.ConsolePrint {
			consolehelper.PrintLogConsole(log.AppName, "error", msg, logs)
		}
		filehelper.PrintLogFile(log.ErrorWriter, log.AppName, "error", msg, logs)
	}
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

func (log *LogHelper) InfoCustom(fields map[string]interface{}) {
	if GetLogLevel(log.LogLevel) <= Log_Info {
		if log.ConsolePrint {
			consolehelper.PrintLogConsoleCustom(log.AppName,"info", fields, logs)
		}
		filehelper.PrintLogFileCustom(log.InfoWriter, log.AppName, "info", fields, logs)
	}
}

func (log *LogHelper) TraceCustom(fields map[string]interface{}) {
	if GetLogLevel(log.LogLevel) <= Log_Trace {
		if log.ConsolePrint {
			consolehelper.PrintLogConsoleCustom(log.AppName,"trace",  fields, logs)
		}
		filehelper.PrintLogFileCustom(log.TraceWriter, log.AppName, "trace", fields, logs)
	}
}

func (log *LogHelper) DebugCustom(fields map[string]interface{}) {
	if GetLogLevel(log.LogLevel) <= Log_Debug {
		if log.ConsolePrint {
			consolehelper.PrintLogConsoleCustom(log.AppName, "debug",fields, logs)
		}
		filehelper.PrintLogFileCustom(log.DebugWriter, log.AppName, "debug", fields, logs)
	}
}

func (log *LogHelper) WarnCustom(fields map[string]interface{}) {
	if GetLogLevel(log.LogLevel) <= Log_Warn {
		if log.ConsolePrint {
			consolehelper.PrintLogConsoleCustom(log.AppName, "warn", fields, logs)
		}
		filehelper.PrintLogFileCustom(log.WarnWriter, log.AppName, "warn", fields, logs)
	}
}

func (log *LogHelper) ErrorCustom(fields map[string]interface{}) {
	if GetLogLevel(log.LogLevel) <= Log_Error {
		if log.ConsolePrint {
			consolehelper.PrintLogConsoleCustom(log.AppName, "error", fields, logs)
		}
		filehelper.PrintLogFileCustom(log.ErrorWriter, log.AppName, "error", fields, logs)
	}
}