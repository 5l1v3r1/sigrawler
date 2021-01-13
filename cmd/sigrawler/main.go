package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/drsigned/gos"
	"github.com/drsigned/sigrawler/pkg/sigrawler"
	"github.com/logrusorgru/aurora/v3"
)

type options struct {
	noColor bool
	silent  bool
	URLs    string
	output  string
}

var (
	co options
	au aurora.Aurora
	so sigrawler.Options
)

func banner() {
	fmt.Fprintln(os.Stderr, aurora.BrightBlue(`
     _                          _
 ___(_) __ _ _ __ __ ___      _| | ___ _ __
/ __| |/ _`+"`"+` | '__/ _`+"`"+` \ \ /\ / / |/ _ \ '__|
\__ \ | (_| | | | (_| |\ V  V /| |  __/ |
|___/_|\__, |_|  \__,_| \_/\_/ |_|\___|_| v1.2.0
       |___/
`).Bold())
}

func init() {
	flag.BoolVar(&so.Debug, "debug", false, "")
	flag.IntVar(&so.RandomDelay, "random-delay", 2, "")
	flag.IntVar(&so.Depth, "depth", 1, "")
	flag.StringVar(&co.URLs, "iL", "", "")
	flag.BoolVar(&so.IncludeSubs, "iS", false, "")
	flag.BoolVar(&co.noColor, "nC", false, "")
	flag.StringVar(&co.output, "oJ", "", "")
	flag.BoolVar(&co.silent, "silent", false, "")
	flag.IntVar(&so.Threads, "threads", 20, "")
	flag.IntVar(&so.Timeout, "timeout", 10, "")
	flag.StringVar(&so.UserAgent, "UA", "", "")
	flag.StringVar(&so.Proxies, "proxies", "", "")

	flag.Usage = func() {
		banner()

		h := "USAGE:\n"
		h += "  sigrawler [OPTIONS]\n"

		h += "\nCRAWLER OPTIONS:\n"
		h += "  -depth           maximum limit on the recursion depth of visited URLs. (default 1)\n"
		h += "  -iS              extend scope to include subdomains (default: false)\n"
		h += "  -proxies         comma separated list of proxies\n"
		h += "  -random-delay    maximum random delay between requests (default: 2s)\n"
		h += "  -threads         maximum no. of concurrent requests (default 20)\n"
		h += "  -timeout         HTTP timeout (default 10s)\n"
		h += "  -UA              User Agent to use\n"

		h += "\nINPUT OPTIONS:\n"
		h += "  -iL              urls to crawl (use `iL -` to read from stdin)\n"

		h += "\nOUTPUT OPTIONS:\n"
		h += "  -debug           stdout: debug mode (default: false)\n"
		h += "  -nC              stdout: no color mode (default: false)\n"
		h += "  -oJ              JSON: output file\n"
		h += "  -silent          stdout: silent mode (default: false)\n"

		fmt.Fprintf(os.Stderr, h)
	}

	flag.Parse()

	au = aurora.NewAurora(!co.noColor)
}

func main() {
	if !co.silent {
		banner()
	}

	options, err := sigrawler.ParseOptions(&so)
	if err != nil {
		log.Fatalln(err)
	}

	URLs := make(chan string)

	go func() {
		defer close(URLs)

		var scanner *bufio.Scanner

		if co.URLs == "-" {
			if !gos.HasStdin() {
				log.Fatalln(errors.New("no stdin: '-iL -' provided"))
			}

			scanner = bufio.NewScanner(os.Stdin)
		} else {
			openedFile, err := os.Open(co.URLs)
			if err != nil {
				log.Fatalln(err)
			}
			defer openedFile.Close()

			scanner = bufio.NewScanner(openedFile)
		}

		for scanner.Scan() {
			if scanner.Text() != "" {
				URLs <- scanner.Text()
			}
		}

		if scanner.Err() != nil {
			log.Fatalln(scanner.Err())
		}
	}()

	var wg sync.WaitGroup
	var output sigrawler.Results

	for URL := range URLs {
		wg.Add(1)

		go func(URL string) {
			defer wg.Done()

			crawler, err := sigrawler.New(URL, options)
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

func saveResults(outputPath string, output sigrawler.Results) error {
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
