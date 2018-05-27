package filehook

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

// VERSION ...
const VERSION = "0.0.2"

// FileHook write log to a file
// LineBreak can choose Windows or Unix-like
// SegmentInterval new log file create interval, second
type FileHook struct {
	File      FileWriter
	formatter log.Formatter
	option    *Option
	mux       sync.Mutex
}

// FileWriter ...
type FileWriter struct {
	mux sync.Mutex
	*os.File
}

// New create the hook
func New(option *Option) (*FileHook, error) {
	hook := new(FileHook)
	parseOption(option)
	hook.option = option

	f := new(FileFormatter)
	f.LineBreak = option.File.LineBreak
	hook.formatter = f

	// run auto create log file routine
	err := hook.createLogFile()
	if err != nil {
		return nil, fmt.Errorf("create log file error: %v", err)
	}
	go hook.fileAutoSegment()

	// run auto compress log files routine
	go hook.autoCompress()

	return hook, nil
}

// SetFormatter set formatter
func (hook *FileHook) SetFormatter(f log.Formatter) {
	hook.mux.Lock()
	defer hook.mux.Unlock()

	hook.formatter = f
}

// Fire writes the log file to defined path.
// User who run this function needs write permissions to the file or directory if the file does not yet exist.
func (hook *FileHook) Fire(entry *log.Entry) error {
	return hook.write(entry)
}

// Levels returns configured log levels.
func (hook *FileHook) Levels() []log.Level {
	return log.AllLevels
}

func (hook *FileHook) write(entry *log.Entry) error {
	serialized, err := hook.formatter.Format(entry)
	if err != nil {
		return fmt.Errorf("file hook format entry error: %v", err)
	}
	hook.File.mux.Lock()
	defer hook.File.mux.Unlock()

	n, err := hook.File.WriteString(string(serialized))
	if err != nil {
		return fmt.Errorf("file hook write %d bytes to file[%v] error: %v", n, hook.File.File.Name(), err)
	}
	return nil
}

func (hook *FileHook) fileAutoSegment() {
	ticker := time.NewTicker(time.Second * time.Duration(hook.option.File.Interval))
	for {
		select {
		case <-ticker.C:
			err := hook.createLogFile()
			if err != nil {
				log.Errorf("create log file error: %v", err)
				return
			}
		}
	}
}

func (hook *FileHook) createLogFile() error {
	_, err := os.Stat(hook.option.Path)
	if err != nil && os.IsNotExist(err) {
		err := os.Mkdir(hook.option.Path, 0777)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	fileName := parseNamePattern(time.Now(), hook.option.NamePattern)
	fileName = generateFileName(filepath.Join(hook.option.Path, fileName), 0)
	fileName += hook.option.File.Ext

	f, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("Can't create file %v, error: %v", fileName, err)
	}
	hook.mux.Lock()
	hook.File = FileWriter{File: f}
	hook.mux.Unlock()
	return nil
}

// parseNamePattern can configure:
// %YY year like  	2018
// %MM nonth like 	03
// %DD date like  	17
// %HH hour like	16(Use 24-hour clock)
// %mm minute like 	03
// %SS second like	02
func parseNamePattern(t time.Time, pattern string) string {
	timeStrings := strings.Split(t.Format("2006-01-02-15-04-05"), "-")
	year := timeStrings[0]
	month := timeStrings[1]
	day := timeStrings[2]
	hour := timeStrings[3]
	minutes := timeStrings[4]
	second := timeStrings[5]

	result := strings.Replace(pattern, "%YY", year, -1)
	result = strings.Replace(result, "%MM", month, -1)
	result = strings.Replace(result, "%DD", day, -1)
	result = strings.Replace(result, "%HH", hour, -1)
	result = strings.Replace(result, "%mm", minutes, -1)
	result = strings.Replace(result, "%SS", second, -1)
	return result
}

func generateFileName(name string, count uint64) string {
	suffix := ""
	if count != 0 {
		suffix = fmt.Sprintf(".%d", count)
	}
	_, err := os.Stat(name + suffix)
	if err == nil {
		return generateFileName(name, count+1)
	}
	return name + suffix
}
