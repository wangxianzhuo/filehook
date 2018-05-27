package main

import (
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/wangxianzhuo/filehook"
)

func main() {
	hookConfig := filehook.NewOption()
	hookConfig.Compress.Enable = true
	hookConfig.Compress.Interval = 10

	log.SetLevel(log.DebugLevel)

	hook, err := filehook.New(hookConfig)
	if err != nil {
		panic(err)
	}

	f := new(filehook.FileFormatter)
	f.TimestampFormat = "2006-01-02 15:04:05.999999999"
	hook.SetFormatter(f)

	log.AddHook(hook)

	log.Infof("info1")
	log.Infof("info2")
	// time.Sleep(15 * time.Second)
	log.Warnf("warn1")
	log.Warnf("warn2")
	time.Sleep(70 * time.Second)
}
