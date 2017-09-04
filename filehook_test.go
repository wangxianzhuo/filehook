package filehook

import (
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
)

func Test_Hook(t *testing.T) {
	hook, err := NewFileHook("", "windows", 10)
	if err != nil {
		t.Fatalf("init FileHook error: %v", err)
	}
	log.AddHook(hook)

	log.Infof("info1")
	log.Infof("info2")
	time.Sleep(15 * time.Second)
	log.Warnf("warn1")
	log.Warnf("warn2")
}
