package crawler

// Options is the structure of the options expected
type Options struct {
	Debug       bool
	Delay       int
	Depth       int
	IncludeSubs bool
	Threads     int
	Proxies     string
	Timeout     int
	UserAgent   string
}

// ParseOptions is a
func ParseOptions(options *Options) (*Options, error) {
	return options, nil
}
