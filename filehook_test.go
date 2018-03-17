package filehook

import (
	"fmt"
	"os"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
)

func Test_Hook(t *testing.T) {
	hook, err := New("", "windows", 0)
	if err != nil {
		t.Fatalf("init FileHook error: %v", err)
	}
	log.AddHook(hook)

	log.Infof("info1")
	log.Infof("info2")
	time.Sleep(15 * time.Second)
	log.Warnf("warn1")
	log.Warnf("warn2")

	for index := 1; index < 20; index++ {
		go func(index int) {
			log.Info(index)
			time.Sleep(500 * time.Millisecond)
		}(index)
	}

	time.Sleep(60 * time.Second)
}

func Test_generateFileName(t *testing.T) {
	name := "./logs/tmp.log"
	files := make([]*os.File, 21)

	defer func() {
		for _, file := range files {
			name := file.Name()
			file.Close()
			err := os.Remove(name)
			if err != nil {
				t.Fatal(err)
			}
		}
	}()

	_, err := os.Stat("./logs")
	if err != nil && os.IsNotExist(err) {
		err := os.Mkdir("./logs", 0777)
		if err != nil {
			t.Fatal(err)
		}
	} else if err != nil {
		t.Fatal(err)
	}

	files[0], err = os.Create(name)
	if err != nil {
		t.Fatal(err)
	}

	if generateFileName(name, 0) != name+".1" {
		t.Fatalf("error name for: %v\n", generateFileName(name, 0))
	}

	for index := 1; index <= 20; index++ {
		files[index], err = os.Create(fmt.Sprintf("%s.%d", name, index))
		if err != nil {
			t.Fatal(err)
		}
	}

	if generateFileName(name, 0) != name+".21" {
		t.Fatalf("error name for: %v\n", generateFileName(name, 0))
	}
}
