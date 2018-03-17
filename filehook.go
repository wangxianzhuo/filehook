package filehook

import (
	"bytes"
	"fmt"
	"os"
	"sort"
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
	SegmentInterval      int64
	Path                 string
	LineBreak            string
	File                 FileWriter
	LatestFileCreateDate time.Time
	namePattern          string
	mux                  sync.Mutex
}

// FileWriter ...
type FileWriter struct {
	mux sync.Mutex
	*os.File
}

// New create the hook
func New(filePath, osType string, segmentInterval int64) (*FileHook, error) {
	hook := new(FileHook)
	if filePath == "" {
		hook.Path = "logs/"
	} else {
		hook.Path = filePath
	}
	if segmentInterval == 0 {
		hook.SegmentInterval = 60 * 60 * 24
	} else {
		hook.SegmentInterval = segmentInterval
	}
	hook.LatestFileCreateDate = time.Now()
	switch osType {
	case "windows":
		hook.LineBreak = "\r\n"
	default:
		hook.LineBreak = "\n"
	}
	if hook.namePattern == "" {
		hook.namePattern = "%YY-%MM-%DD_%HH-%mm-%SS.log"
	}

	err := hook.fileAutoSegment(false)
	if err != nil {
		return nil, err
	}
	go hook.fileAutoSegment(true)

	return hook, nil
}

// Fire writes the log file to defined path.
// User who run this function needs write permissions to the file or directory if the file does not yet exist.
func (hook *FileHook) Fire(entry *log.Entry) error {
	return hook.writeLog(entry)
}

// Levels returns configured log levels.
func (hook *FileHook) Levels() []log.Level {
	return log.AllLevels
}

func (hook *FileHook) writeLog(entry *log.Entry) error {
	line, err := entry.String()
	if err != nil {
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

	hook.File.mux.Lock()
	defer hook.File.mux.Unlock()

	_, err = hook.File.Write([]byte(line))
	if err != nil {
		return err
	}
	return nil
}

func (hook *FileHook) fileAutoSegment(isInited bool) error {
	if isInited && hook.SegmentInterval > 0 {
		for {
			now := time.Now()
			recreateDate := hook.LatestFileCreateDate.Add(time.Duration(hook.SegmentInterval) * time.Second)
			if now.Before(recreateDate) {
				time.Sleep(recreateDate.Sub(now))
			}
			hook.createLogFile()
		}
	} else if !isInited {
		return hook.createLogFile()
	}
	return nil
}

func (hook *FileHook) createLogFile() error {
	hook.LatestFileCreateDate = time.Now()
	_, err := os.Stat(hook.Path)
	if err != nil && os.IsNotExist(err) {
		err := os.Mkdir(hook.Path, 0777)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	fileName := parseNamePattern(hook.LatestFileCreateDate, hook.namePattern)
	fileName = generateFileName(hook.Path+fileName, 0)

	f, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("Can't create a file hook, do not use log files, error: %v", err)
	}
	hook.File.Close()
	hook.mux.Lock()
	hook.File = FileWriter{File: f}
	hook.mux.Unlock()
	return nil
}

// SetFileNamePattern ...
func (hook *FileHook) SetFileNamePattern(newPattern string) {
	hook.mux.Lock()
	defer hook.mux.Unlock()

	hook.namePattern = newPattern
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
