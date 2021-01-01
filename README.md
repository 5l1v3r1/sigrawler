# sigrawler

[![release](https://img.shields.io/github/release/drsigned/sigrawler?style=flat&color=0040ff)](https://github.com/drsigned/sigrawler/releases) [![maintenance](https://img.shields.io/badge/maintained%3F-yes-0040ff.svg)](https://github.com/drsigned/sigrawler) [![open issues](https://img.shields.io/github/issues-raw/drsigned/sigrawler.svg?style=flat&color=0040ff)](https://github.com/drsigned/sigrawler/issues?q=is:issue+is:open) [![closed issues](https://img.shields.io/github/issues-closed-raw/drsigned/sigrawler.svg?style=flat&color=0040ff)](https://github.com/drsigned/sigrawler/issues?q=is:issue+is:closed) [![license](https://img.shields.io/badge/license-MIT-gray.svg?colorB=0040FF)](https://github.com/drsigned/sigrawler/blob/master/LICENSE) [![twitter](https://img.shields.io/badge/twitter-@drsigned-0040ff.svg)](https://twitter.com/drsigned)

## Resources

* [Usage](#usage)
* [Installation](#installation)
    * [From Binary](#from-binary)
    * [From source](#from-source)
    * [From github](#from-github)
* [Contribution](#contribution)

## Usage

```text
$ sigrawler -h

     _                          _
 ___(_) __ _ _ __ __ ___      _| | ___ _ __
/ __| |/ _` | '__/ _` \ \ /\ / / |/ _ \ '__|
\__ \ | (_| | | | (_| |\ V  V /| |  __/ |
|___/_|\__, |_|  \__,_| \_/\_/ |_|\___|_| v1.1.0
       |___/

USAGE:
  sigrawler [OPTIONS]

OPTIONS:
  -debug          debug mode (default: false)
  -delay          delay between requests. (default 5s)
  -depth          maximum limit on the recursion depth of visited URLs. (default 1)
  -iL             urls to crawl (use `iL -` to read from stdin)
  -iS             extend scope to include subdomains (default: false)
  -nC             no color mode
  -oJ             JSON output file
  -s              silent mode: print urls only (default: false)
  -threads        maximum no. of concurrent requests (default 20)
  -timeout        HTTP timeout (default 10s)
  -UA             User Agent to use
  -x              comma separated list of proxies
```

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

## Contibution

[Issues](https://github.com/drsigned/sigrawler/issues) and [Pull Requests](https://github.com/drsigned/sigrawler/pulls) are welcome! 