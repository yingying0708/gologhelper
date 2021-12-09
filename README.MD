## 简介
此项目是一个开源的golang的日志功能项目

## 核心技术及组件
logrus v1.8.1
file-rotatelogs v2.4.0+incompatible

## 获取
go get github.com/yingying0708/gologhelper

## 使用
log := gologhelper.GetLogHelper("项目名称", "日志路径").SetConsolePrint(true).SetLogLevel("info").SetWhen("m").SetBackupCount(2)

log.Info("info")

## 方法介绍
SetConsolePrint(true)  控制台打印，默认是false


SetLogLevel("info")    只生成info的日志文件，默认是error，日志级别从小到大排列是：track，debug，info，warn，error。例：设置info ，会打印info、warn、error的类型的日志


SetWhen("m")           设置按什么进行分割，默认是D，D：天，H：小时，m：分钟


SetBackupCount(2)      保留多少个文件，默认是15

