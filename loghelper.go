package loghelper

import (
	consolehelper "github.com/yingying0708/loghelper/ConsoleLogPrint"
	filehelper "github.com/yingying0708/loghelper/FileLogPrint"
	"strings"
	"time"

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
	Writer       *rotatelogs.RotateLogs
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
		Writer:       getWriter(logPath, "D", 15),
	}
}

//设置when(D:天，H：小时，M：分钟，默认是D)
func (log *LogHelper) SetWhen(when string) *LogHelper {
	if when != "" {
		log.When = strings.ToUpper(when)
		logPath := log.LogPath + log.AppName + "_p1_" + log.LogLevel + ".log"
		log.Writer = getWriter(logPath, log.When, log.BackupCount)
	}
	return log
}

//设置日志级别（error,debug,info,trace,warn 默认是error）
func (log *LogHelper) SetLogLevel(level string) *LogHelper {
	levelstr := "error"
	if level != "" {
		levelstr = strings.ToLower(level)
	}
	if levelstr == "debug" {
		logs.SetLevel(logrus.DebugLevel)
	}
	if levelstr == "info" {
		logs.SetLevel(logrus.InfoLevel)
	}
	if levelstr == "trace" {
		logs.SetLevel(logrus.TraceLevel)
	}
	if levelstr == "warn" {
		logs.SetLevel(logrus.WarnLevel)
	}
	if levelstr == "error" {
		logs.SetLevel(logrus.ErrorLevel)
	}
	log.LogLevel = levelstr
	logPath := log.LogPath + log.AppName + "_p1_" + log.LogLevel + ".log"
	log.Writer = getWriter(logPath, log.When, log.BackupCount)
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
		logPath := log.LogPath + log.AppName + "_p1_" + log.LogLevel + ".log"
		log.Writer = getWriter(logPath, log.When, log.BackupCount)
	}
	return log
}

func (log *LogHelper) Info(msg interface{}) {
	if log.LogLevel == "info" {
		if log.ConsolePrint {
			consolehelper.PrintLogConsole(log.AppName, msg, logs)
		}
		filehelper.PrintLogFile(log.Writer, log.AppName, msg, logs)
	}
}

func (log *LogHelper) Trace(msg interface{}) {
	if log.LogLevel == "trace" {
		if log.ConsolePrint {
			consolehelper.PrintLogConsole(log.AppName, msg, logs)
		}
		filehelper.PrintLogFile(log.Writer, log.AppName, msg, logs)
	}
}

func (log *LogHelper) Debug(msg interface{}) {
	if log.LogLevel == "debug" {
		if log.ConsolePrint {
			consolehelper.PrintLogConsole(log.AppName, msg, logs)
		}
		filehelper.PrintLogFile(log.Writer, log.AppName, msg, logs)
	}
}

func (log *LogHelper) Warn(msg interface{}) {
	if log.LogLevel == "warn" {
		if log.ConsolePrint {
			consolehelper.PrintLogConsole(log.AppName, msg, logs)
		}
		filehelper.PrintLogFile(log.Writer, log.AppName, msg, logs)
	}
}

func (log *LogHelper) Error(msg interface{}) {
	if log.LogLevel == "error" {
		if log.ConsolePrint {
			consolehelper.PrintLogConsole(log.AppName, msg, logs)
		}
		filehelper.PrintLogFile(log.Writer, log.AppName, msg, logs)
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
