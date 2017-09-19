package filehook

import (
	"bytes"
	"fmt"
	"os"
	"sort"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

const VERSION = "0.0.1"

// FileHook write log to a file
// LineBreak can choose Windows or Unix-like
// RecreateFileInterval new log file create interval, second
type FileHook struct {
	RecreateFileInterval int64
	Path                 string
	LineBreak            string
	File                 FileWriter
	LatestFileCreateDate time.Time
	Mux                  sync.Mutex
}

type FileWriter struct {
	Mux sync.Mutex
	*os.File
}

func New(filePath, osType string, recreateFileInterval int64) (*FileHook, error) {
	hook := new(FileHook)

	if filePath == "" {
		hook.Path = "logs/"
	} else {
		hook.Path = filePath
	}

	if recreateFileInterval <= 0 {
		hook.RecreateFileInterval = 60 * 60 * 24
	} else {
		hook.RecreateFileInterval = recreateFileInterval
	}

	switch osType {
	case "windows":
		hook.LineBreak = "\r\n"
	default:
		hook.LineBreak = "\n"
	}

	err := hook.createLogFile()
	if err != nil {
		return nil, err
	}

	return hook, nil
}
func (hook *FileHook) Fire(entry *log.Entry) error {
	err := hook.createLogFile()
	if err != nil {
		return err
	}
	err = hook.writeLog(entry)
	if err != nil {
		return err
	}
	return nil
}
func (hook *FileHook) Levels() []log.Level {
	return log.AllLevels
}
func (hook *FileHook) UseNewFile() error {
	err := hook.createLogFile()
	if err != nil {
		return err
	}
	return nil
}

func (hook *FileHook) writeLog(entry *log.Entry) error {
	line, err := entry.String()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read entry, %v", err)
		return err
	}
	keys := make([]string, 0, len(entry.Data))
	for k := range entry.Data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	buffer := bytes.Buffer{}
	for _, k := range keys {
		buffer.WriteString(k)
		buffer.WriteString("=")

		switch value := entry.Data[k].(type) {
		case string:
			buffer.WriteString(value)
		case error:
			errmsg := value.Error()
			buffer.WriteString(errmsg)
		default:
			fmt.Fprint(&buffer, value)
		}
		buffer.WriteString("\t")
	}
	var level string
	switch entry.Level {
	case log.PanicLevel:
		level = "PANI"
	case log.FatalLevel:
		level = "FATA"
	case log.ErrorLevel:
		level = "ERRO"
	case log.WarnLevel:
		level = "WARN"
	case log.InfoLevel:
		level = "INFO"
	case log.DebugLevel:
		level = "DEBU"
	}
	line = fmt.Sprintf("%s[%v] %-80s\t%s"+hook.LineBreak, level, entry.Time.Format("2006-01-02 15:04:05"), entry.Message, buffer.String())

	hook.File.Mux.Lock()
	defer hook.File.Mux.Unlock()
	_, err = hook.File.Write([]byte(line))
	if err != nil {
		return err
	}
	return nil
}

func (hook *FileHook) createLogFile() error {
	now := time.Now()
	recreateDate := hook.LatestFileCreateDate.Add(time.Duration(hook.RecreateFileInterval) * time.Second)

	// fmt.Printf("now: %v\trecreate date: %v\n", now.Format("2006-01-02 15:04:05"), recreateDate.Format("2006-01-02 15:04:05"))

	if now.Before(recreateDate) {
		return nil
	}

	_, err := os.Stat(hook.Path)
	if err != nil {
		err := os.Mkdir(hook.Path, 0777)
		if err != nil {
			return err
		}
	}

	f, err := os.Create(hook.Path + "log." + now.Format("2006-01-02_15-04-05"))
	if err != nil {
		return fmt.Errorf("Can't create a file hook, do not use log files, error: %v", err)
	}
	hook.File.Close()
	hook.Mux.Lock()
	hook.File = FileWriter{File: f}
	hook.Mux.Unlock()

	hook.LatestFileCreateDate = now
	return nil
}
