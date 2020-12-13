package crawler

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/drsigned/gos"
	ext "github.com/drsigned/sigrawler/pkg/crawler/extensions"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/debug"
	"github.com/gocolly/colly/v2/extensions"
)

// Crawler is a
type Crawler struct {
	Options    *Options
	URL        *gos.URL
	PCollector *colly.Collector
	JCollector *colly.Collector
}

// Results is a
type Results struct {
	URLs    []string `json:"urls,omitempty"`
	Buckets []string `json:"s3,omitempty"`
}

// New is a
func New(URL string, options *Options) (crawler Crawler, err error) {
	crawler.Options = options

	parsedURL, err := gos.ParseURL(URL)
	if err != nil {
		return crawler, err
	}

	crawler.URL = parsedURL

	eTLDPlus1 := parsedURL.ETLDPlus1
	escapedETLDPlus1 := strings.ReplaceAll(eTLDPlus1, ".", "\\.")

	pCollector := colly.NewCollector(
		colly.MaxDepth(options.Depth),
	)

	if options.IncludeSubs {
		pCollector.URLFilters = []*regexp.Regexp{
			regexp.MustCompile(
				fmt.Sprintf(`(https?)://[^\s?#\/]*%s/?[^\s]*`, escapedETLDPlus1),
			),
		}

	} else {
		pCollector.AllowedDomains = []string{
			eTLDPlus1,
			"www." + eTLDPlus1,
		}
	}

	// Set User-Agent
	if options.UserAgent != "" {
		pCollector.UserAgent = options.UserAgent
	} else {
		extensions.RandomMobileUserAgent(pCollector)
	}

	// Referer
	extensions.Referer(pCollector)

	// Debug
	if options.Debug {
		pCollector.SetDebugger(&debug.LogDebugger{})
	}

	// Insecure
	if options.Insecure {
		pCollector.WithTransport(&http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		})
	}

	pCollector.SetRequestTimeout(
		time.Duration(crawler.Options.Timeout) * time.Second,
	)

	err = pCollector.Limit(&colly.LimitRule{
		DomainGlob:  fmt.Sprintf("*%s", parsedURL.ETLDPlus1),
		Parallelism: crawler.Options.Concurrency,
		Delay:       time.Duration(crawler.Options.Delay) * time.Second,
		RandomDelay: time.Duration(crawler.Options.RandomDelay) * time.Second,
	})
	if err != nil {
		return crawler, err
	}

	jCollector := pCollector.Clone()
	jCollector.URLFilters = nil

	crawler.PCollector = pCollector
	crawler.JCollector = jCollector

	return crawler, nil
}

// Run is a
func (crawler *Crawler) Run(URL string) (results Results, err error) {
	// make sure the url has been set
	if URL == "" {
		return results, errors.New("no url was provided")
	}

	// these will store the discovered assets to avoid duplicates
	var URLs sync.Map
	URLsSlice := make([]string, 0)
	var buckets sync.Map
	bucketsSlice := make([]string, 0)

	jsRegex := regexp.MustCompile(`(?m).*?\.*(js|json)(\?.*?|)$`)
	filesRegex := regexp.MustCompile(`(?m).*?\.*(jpg|png|gif|webp|tiff|psd|raw|bmp|heif|ico|css|pdf|jpeg|css|tif|ttf|woff|woff2|pdf|doc|svg|txt|mp3|mp4|eot)(\?.*?|)$`)

	crawler.PCollector.OnRequest(func(request *colly.Request) {
		reqURL := request.URL.String()

		// JavaScript
		if match := jsRegex.MatchString(reqURL); match {
			// Minified JavaScript
			if strings.Contains(reqURL, ".min.js") {
				js := strings.ReplaceAll(reqURL, ".min.js", ".js")

				if _, exists := URLs.Load(js); exists {
					return
				}

				crawler.JCollector.Visit(js)

				URLs.Store(js, struct{}{})
			}

			crawler.JCollector.Visit(reqURL)

			// Cancel the request to ensure we don't process it on this collector
			request.Abort()
			return
		}

		// Files: Is it an image or similar? Don't request it.
		if match := filesRegex.MatchString(reqURL); match {
			request.Abort()
			return
		}
	})

	crawler.PCollector.OnResponse(func(response *colly.Response) {
		// s3 buckets
		S3s, err := ext.S3finder(string(response.Body))
		if err != nil {
			return
		}

		for _, S3 := range S3s {
			if _, exists := buckets.Load(S3); exists {
				return
			}

			fmt.Println("[s3]", S3)

			bucketsSlice = append(bucketsSlice, S3)
			buckets.Store(S3, struct{}{})
		}
	})

	crawler.PCollector.OnHTML("[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// Get the absolute URL
		absoluteURL := e.Request.AbsoluteURL(link)
		// Trim the trailing slash
		absoluteURL = strings.TrimRight(absoluteURL, "/")
		// Trim the spaces on either end (if any)
		absoluteURL = strings.Trim(absoluteURL, " ")

		URL := fixURL(absoluteURL, crawler.URL)
		if URL == "" {
			return
		}

		if _, exists := URLs.Load(URL); exists {
			return
		}

		e.Request.Visit(link)

		if ok := crawler.record("[url]", URL); ok {
			URLsSlice = append(URLsSlice, URL)
		}

		URLs.Store(URL, struct{}{})
	})

	crawler.PCollector.OnHTML("[src]", func(e *colly.HTMLElement) {
		link := e.Attr("src")
		// Get the absolute URL
		absoluteURL := e.Request.AbsoluteURL(link)

		URL := fixURL(absoluteURL, crawler.URL)
		if URL == "" {
			return
		}

		crawler.PCollector.Visit(URL)

		crawler.record("[javascript]", URL)

		URLs.Store(URL, struct{}{})
	})

	crawler.JCollector.OnResponse(func(response *colly.Response) {
		endpoints, err := ext.Linkfinder(string(response.Body))
		if err != nil {
			return
		}

		for _, endpoint := range endpoints {
			// Skip blank entries
			if len(endpoint) <= 0 {
				continue
			}
			// Remove the single and double quotes from the parsed link on the ends
			endpoint = strings.Trim(endpoint, "\"")
			endpoint = strings.Trim(endpoint, "'")
			// Get the absolute URL
			absoluteURL := response.Request.AbsoluteURL(endpoint)

			if _, exists := URLs.Load(absoluteURL); exists {
				return
			}

			crawler.PCollector.Visit(absoluteURL)

			if ok := crawler.record("[linkfinder]", absoluteURL); ok {
				URLsSlice = append(URLsSlice, absoluteURL)
			}

			URLs.Store(absoluteURL, struct{}{})
		}

		// s3 buckets
		S3s, err := ext.S3finder(string(response.Body))
		if err != nil {
			return
		}

		for _, S3 := range S3s {
			if _, exists := buckets.Load(S3); exists {
				return
			}

			fmt.Println("[s3]", S3)

			bucketsSlice = append(bucketsSlice, S3)
			buckets.Store(S3, struct{}{})
		}
	})

	// setup a waitgroup to run all methods at the same time
	var wg sync.WaitGroup

	// colly
	wg.Add(1)
	go func() {
		defer wg.Done()

		crawler.PCollector.Visit(crawler.URL.String())
	}()

	wg.Wait()

	results.URLs = URLsSlice
	results.Buckets = bucketsSlice

	return results, nil
}

func (crawler *Crawler) record(tag string, URL string) (print bool) {
	parsedURL, err := gos.ParseURL(URL)
	if err != nil {
		return false
	}

	if crawler.Options.IncludeSubs {
		escapedHost := strings.ReplaceAll(crawler.URL.Host, ".", "\\.")
		print, _ = regexp.MatchString(".*(\\.|\\/\\/)"+escapedHost+"((#|\\/|\\?).*)?", URL)
	} else {
		print = parsedURL.Host == crawler.URL.Host || parsedURL.Host == "www."+crawler.URL.Host
	}

	if print {
		fmt.Println(tag, URL)
	}

	return print
}
