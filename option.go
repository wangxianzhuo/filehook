package filehook

// Option contains all configures for this hook
type Option struct {
	Path        string
	NamePattern string
	File        FileConf
	Compress    CompressConf
}
type FileConf struct {
	Interval  int64
	LineBreak string
	Ext       string
}
type CompressConf struct {
	Enable   bool
	Interval int64
	Ext      string
}

func NewOption() *Option {
	option := Option{}

	option.Path = "logs"
	option.NamePattern = "%YY-%MM-%DD_%HH-%mm-%SS"

	option.File.Interval = 86400
	option.File.LineBreak = "\n"
	option.File.Ext = ".log"

	option.Compress.Enable = false
	option.Compress.Interval = option.File.Interval * 30
	option.Compress.Ext = ".tar.gz"

	return &option
}

func parseOption(option *Option) {
	if option == nil {
		option = NewOption()
	}

	if option.Path == "" {
		option.Path = "logs"
	}
	if option.NamePattern == "" {
		option.NamePattern = "%YY-%MM-%DD_%HH-%mm-%SS"
	}

	if option.File.Ext == "" {
		option.File.Ext = ".log"
	}
	if option.File.Interval <= 0 {
		option.File.Interval = 86400
	}
	if option.File.LineBreak == "" {
		option.File.LineBreak = "\n"
	}

	if option.Compress.Interval <= 0 {
		option.Compress.Interval = option.File.Interval * 30
	}
	if option.Compress.Ext == "" {
		option.Compress.Ext = ".tar.gz"
	}
}
