package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/drsigned/gos"
	"github.com/drsigned/sigrawler/pkg/crawler"
	"github.com/logrusorgru/aurora/v3"
)

type options struct {
	noColor bool
	silent  bool
	URL     string
	output  string
}

var (
	co options
	au aurora.Aurora
	so crawler.Options
)

func banner() {
	fmt.Fprintln(os.Stderr, aurora.BrightBlue(`
     _                          _
 ___(_) __ _ _ __ __ ___      _| | ___ _ __
/ __| |/ _`+"`"+` | '__/ _`+"`"+` \ \ /\ / / |/ _ \ '__|
\__ \ | (_| | | | (_| |\ V  V /| |  __/ |
|___/_|\__, |_|  \__,_| \_/\_/ |_|\___|_| v1.1.0
       |___/
`).Bold())
}

func init() {
	// GENERAL OPTIONS
	flag.BoolVar(&co.noColor, "nc", false, "")
	flag.StringVar(&co.URL, "u", "", "")

	// CRAWLER OPTIONS
	flag.BoolVar(&so.Debug, "debug", false, "")
	flag.IntVar(&so.Depth, "depth", 1, "")
	flag.IntVar(&so.Delay, "delay", 2000, "")
	flag.IntVar(&so.Threads, "threads", 20, "")
	flag.BoolVar(&so.IncludeSubs, "subs", false, "")
	flag.StringVar(&so.Proxies, "x", "", "")
	flag.IntVar(&so.Timeout, "timeout", 10, "")

	flag.StringVar(&so.UserAgent, "UA", "", "")

	// OUTPUT OPTIONS
	flag.StringVar(&co.output, "o", "", "")
	flag.BoolVar(&co.silent, "s", false, "")

	flag.Usage = func() {
		banner()

		h := "USAGE:\n"
		h += "  sigrawler [OPTIONS]\n"

		h += "\nGENERAL OPTIONS:\n"
		h += "  -nc                no color mode\n"
		h += "  -u                 the url that you wish to crawl\n"

		h += "\nCRAWLER OPTIONS:\n"
		h += "  -debug             debug mode: extremely verbose output (default: false)\n"
		h += "  -delay             delay in ms between requests. (default 2000)\n"
		h += "  -depth             maximum depth to crawl (default: 1)\n"
		h += "  -threads           maximum no. of concurrent requests (default 20)\n"
		h += "  -timeout           HTTP timeout\n"
		h += "  -subs              crawl subdomains (default: false)\n"
		h += "  -UA                User Agent to use\n"
		h += "  -x                 comma separated list of proxies\n"

		h += "\nOUTPUT OPTIONS:\n"
		h += "  -o                 JSON output file\n"
		h += "  -s                 silent mode: print urls only (default: false)\n"

		fmt.Fprintf(os.Stderr, h)
	}

	flag.Parse()

	au = aurora.NewAurora(!co.noColor)
}

func main() {
	options, err := crawler.ParseOptions(&so)
	if err != nil {
		log.Fatalln(err)
	}

	if !co.silent {
		banner()
	}

	URLs := make(chan string, 1)

	if co.URL != "" {
		URLs <- co.URL

		close(URLs)
	} else {
		if !gos.HasStdin() {
			os.Exit(1)
		}

		go func() {
			defer close(URLs)

			scanner := bufio.NewScanner(os.Stdin)

			for scanner.Scan() {
				URL := strings.ToLower(scanner.Text())

				if URL != "" {
					URLs <- URL
				}
			}
		}()
	}

	var wg sync.WaitGroup
	var output crawler.Results

	for URL := range URLs {
		wg.Add(1)

		if URL == "" {
			continue
		}

		go func(URL string) {
			defer wg.Done()

			crawler, err := crawler.New(URL, options)
			if err != nil {
				log.Fatalln(err)
			}

			results, err := crawler.Run(URL)
			if err != nil {
				log.Fatalln(err)
			}

			output.URLs = append(output.URLs, results.URLs...)
			output.Buckets = append(output.Buckets, results.Buckets...)
		}(URL)
	}

	wg.Wait()

	if co.output != "" {
		if err := saveResults(co.output, output); err != nil {
			log.Fatalln(err)
		}
	}
}

func saveResults(outputPath string, output crawler.Results) error {
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		directory, filename := path.Split(outputPath)

		if _, err := os.Stat(directory); os.IsNotExist(err) {
			if directory != "" {
				err = os.MkdirAll(directory, os.ModePerm)
				if err != nil {
					return err
				}
			}
		}

		if strings.ToLower(path.Ext(filename)) != ".json" {
			outputPath = outputPath + ".json"
		}
	}

	outputJSON, err := json.MarshalIndent(output, "", "\t")
	if err != nil {
		return err
	}

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}

	defer outputFile.Close()

	_, err = outputFile.WriteString(string(outputJSON))
	if err != nil {
		return err
	}

	return nil
}
