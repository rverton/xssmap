# xssmap

xssmap is a tiny tool to scan for (DOM)XSS vulnerabilities by using a headless browser based on webkit (PhantomJS). This enables evaluation of attack vectors and is nearly false positive free.

## Requirements

* [PhantomJS (command line tool)](http://phantomjs.org/download.html)

## Installation

1. Install the phantomjs command line tool.
2. `go get github.com/rverton/xssmap`
3. `xssmap -h`

## Usage

    XSSMAP v0.1. github.com/rverton/xssmap

    Usage:
        xssmap [--method=<method>] [--data=<data>] [--json] [--failed] [--payloads=<payloads.txt>] URL
        xssmap -h | --help
        xssmap --version

    Arguments:
        URL  insert {XSS} as a placeholder for payloads

    Options:
        --failed  Show failed attempts.
        --json  Use JSON as output format.
        -h --help  Show this screen.
        --version  Show version.

    Example:
        xssmap http://server.com/foo{XSS}
        xssmap --method=POST --data="foo={XSS}" http://server.com/vuln
        xssmap --failed http://server.com/foo#{XSS}

## Payloads

Payloads are located in payloads.txt. xssmap checks if `window.xssmap` is defined. All payloads are either

* setting `window.xssmap = true;` or
* calling `xssm()`, a function which is injected and executes `window.xssmap = true`.

Most of payloads were slightly modified from [ra2-dom-xss-scanner](https://code.google.com/p/ra2-dom-xss-scanner/).

## License

MIT
