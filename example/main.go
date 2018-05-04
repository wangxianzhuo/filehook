package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/wangxianzhuo/filehook"
)

func main() {
	// use default path(./log)
	hook, err := filehook.New(&filehook.Option{
		Path:            "./logs/",
		SegmentInterval: 86400,
		NamePattern:     "%YY-%MM-%DD_%HH-%mm-%SS.log",
		LineBreak:       "\n",
	})
	if err != nil {
		panic(err)
	}
	log.AddHook(hook)

	log.Infof("info1")
	log.Infof("info2")
	// time.Sleep(15 * time.Second)
	log.Warnf("warn1")
	log.Warnf("warn2")
}
