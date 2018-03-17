# File Hook

File hook for logrus

## Usage

```golang
package main

import (
	logrus "github.com/sirupsen/logrus"
	"github.com/wangxianzhuo/filehook"
	"os"
)

func main() {

	logrus.SetFormatter(&logrus.JSONFormatter{})

	logrus.SetOutput(os.Stderr)

	logrus.SetLevel(logrus.DebugLevel)

	hook, err := filehook.New("", "windows", 0)
	if err != nil {
		panic(err)
	}
	hook.SetFileNamePattern("log%YY%MM%DD")
	logrus.AddHook(hook)

	logrus.Warn("warn")
	logrus.Info("info")
	logrus.Debug("debug")
}
```

## Parameter

Optional

- Path `the logs store path`
	- type `string`
	- default `default ./logs/`
- LineBreak `the line break`
	- type `string`
	- value `\n` or `\r\n`
	- default `default \n`
- SegmentInterval `segment file interval, unit second`
	- type `int64`
	- value `-1` is no segmention, `0` default segmention, `>0` segmention interval 
	- default `default 86400`
- namePattern is the log file name pattern. It can only use SetFileNamePattern() to modify.
	- type `string`
	- default `%YY-%MM-%DD_%HH-%mm-%SS.log`

## Installation

```shell
go get -u github.com/wangxianzhuo/filehook
```