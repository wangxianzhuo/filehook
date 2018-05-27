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
	log.AddHook(hook)

	log.Infof("info1")
	log.Infof("info2")
	// time.Sleep(15 * time.Second)
	log.Warnf("warn1")
	log.Warnf("warn2")
	time.Sleep(70 * time.Second)
}
