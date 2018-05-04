package filehook

// Option contains all configures for this hook
type Option struct {
	Path            string
	SegmentInterval int64
	NamePattern     string
	LineBreak       string
}
