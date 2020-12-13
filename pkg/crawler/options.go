package crawler

// Options is the structure of the options expected
type Options struct {
	Concurrency int
	Debug       bool
	Delay       int
	Depth       int
	IncludeSubs bool
	Insecure    bool
	RandomDelay int
	Timeout     int
	UserAgent   string
}

// ParseOptions is a
func ParseOptions(options *Options) (*Options, error) {
	return options, nil
}
