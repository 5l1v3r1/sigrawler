# sigrawler

![made with go](https://img.shields.io/badge/made%20with-Go-0040ff.svg) ![maintenance](https://img.shields.io/badge/maintained%3F-yes-0040ff.svg) [![open issues](https://img.shields.io/github/issues-raw/drsigned/sigrawler.svg?style=flat&color=0040ff)](https://github.com/drsigned/sigrawler/issues?q=is:issue+is:open) [![closed issues](https://img.shields.io/github/issues-closed-raw/drsigned/sigrawler.svg?style=flat&color=0040ff)](https://github.com/drsigned/sigrawler/issues?q=is:issue+is:closed) [![license](https://img.shields.io/badge/License-MIT-gray.svg?colorB=0040FF)](https://github.com/drsigned/sigrawler/blob/master/LICENSE.md) [![twitter](https://img.shields.io/badge/twitter-@drsigned-0040ff.svg)](https://twitter.com/drsigned)

## Resources

* [Installation](#installation)
    * [From Binary](#from-binary)
    * [From source](#from-source)
    * [From github](#from-github)
* [Usage](#usage)
* [Contribution](#contribution)

## Installation

#### From Binary

You can download the pre-built binary for your platform from this repository's [releases](https://github.com/drsigned/sigrawler/releases/) page, extract, then move it to your `$PATH`and you're ready to go.

#### From Source

sigrawler requires **go1.14+** to install successfully. Run the following command to get the repo

```bash
$ GO111MODULE=on go get -u -v github.com/drsigned/sigrawler/cmd/sigrawler
```

#### From Github

```bash
$ git clone https://github.com/drsigned/sigrawler.git; cd sigrawler/cmd/sigrawler/; go build; mv sigrawler /usr/local/bin/; sigrawler -h
```

## Usage

```text
$ sigrawler -h

     _                          _
 ___(_) __ _ _ __ __ ___      _| | ___ _ __
/ __| |/ _` | '__/ _` \ \ /\ / / |/ _ \ '__|
\__ \ | (_| | | | (_| |\ V  V /| |  __/ |
|___/_|\__, |_|  \__,_| \_/\_/ |_|\___|_| v1.0.0
       |___/

USAGE:
  sigrawler [OPTIONS]

GENERAL OPTIONS:
  -debug             debug mode: extremely verbose output (default: false)
  -nc                no color mode

CRAWLER OPTIONS:
  -c                 number of maximum allowed concurrent requests of the matching domains (default 5)
  -depth             maximum depth to crawl (default: 1)
  -subs              crawl subdomains (default: false)

HTTP OPTIONS:
  -insecure          ignore invalid HTTPS certificates
  -timeout           HTTP timeout
  -url               the url that you wish to crawl
  -UA                User Agent to use

OUTPUT OPTIONS:
  -o                 JSON output file
  -s                 silent mode: print urls only (default: false)
```

## Contibution

[Issues](https://github.com/drsigned/sigrawler/issues) and [Pull Requests](https://github.com/drsigned/sigrawler/pulls) are welcome! 