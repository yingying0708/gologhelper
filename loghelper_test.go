package loghelper_test

import (
	"testing"
	"time"

	"github.com/yingying0708/loghelper"
)

func TestInfo(t *testing.T) {
	log := loghelper.GetLogHelper("a", "C:\\logs\\").SetConsolePrint(true).SetLogLevel("info").SetWhen("m").SetBackupCount(2)
	for i := 0; i < 200; i++ {
		time.Sleep(time.Second)
		log.Info(i)
	}
}
