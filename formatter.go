package filehook

import (
	"bytes"
	"fmt"
	"sort"
	"time"

	log "github.com/sirupsen/logrus"
)

const defaultTimestampFormat = time.RFC3339

var (
	baseTimestamp time.Time
)

func init() {
	baseTimestamp = time.Now()
}

// FileFormatter ...
type FileFormatter struct {
	// TimestampFormat sets the format used for marshaling timestamps.
	TimestampFormat string

	// LineBreak \n or \r\n
	LineBreak string
}

// Format ...
func (f *FileFormatter) Format(entry *log.Entry) ([]byte, error) {
	var buffer *bytes.Buffer
	keys := make([]string, 0, len(entry.Data))
	for k := range entry.Data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	if entry.Buffer != nil {
		buffer = entry.Buffer
	} else {
		buffer = &bytes.Buffer{}
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

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = defaultTimestampFormat
	}
	lineBreak := f.LineBreak
	if lineBreak == "" {
		lineBreak = "\n"
	}

	fmt.Fprintf(buffer, "%s[%v]\t%-80s\t", level, entry.Time.Format(timestampFormat), entry.Message)

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
			fmt.Fprint(buffer, value)
		}
		buffer.WriteString("\t")
	}

	buffer.WriteString(lineBreak)
	return buffer.Bytes(), nil
}
