package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/docopt/docopt-go"
)

type Request struct {
	url    string
	method string
	data   string
}

type Result struct {
	Vector  XSSVector
	Request Request
	Status  bool
}

type XSSVector struct {
	Payload string
}

func main() {

	usage := `XSSMAP v0.1. github.com/rverton/xssmap

Usage:
	xssmap [--method=<method>] [--data=<data>] [--json] [--failed] [--payloads=<payloads.txt>] URL
	xssmap -h | --help
	xssmap --version

Arguments:
	URL  insert {XSS} as a placeholder for payloads

Options:
	--failed  Also return failed attempts.
	--json  Use JSON as output format.
	-h --help  Show this screen.
	--version  Show version.

Example:
	xssmap http://server.com/foo{XSS}
	xssmap --method=POST --data="foo={XSS}" http://server.com/foo
`

	var outputFormat string = "plain"
	var payloadsFile string = "payloads.txt"
	var showFailed = false
	var success int = 0

	arguments, _ := docopt.Parse(usage, nil, true, "XSSMAP 0.1", false)

	req := Request{
		url:    arguments["URL"].(string),
		method: "GET",
	}

	if val, ok := arguments["--method"]; ok && val != nil {
		req.method = val.(string)
	}

	if val, ok := arguments["--data"]; ok && val != nil {
		req.data = val.(string)
	}

	if val, ok := arguments["--json"]; ok && val != false {
		outputFormat = "json"
	}

	if val, ok := arguments["--failed"]; ok && val != false {
		showFailed = true
	}

	if val, ok := arguments["--payloads"]; ok && val != nil {
		payloadsFile = val.(string)
	}

	results := make(chan Result)

	file, err := os.Open(payloadsFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	go Run(file, req, results)

	for result := range results {
		status := "Failed"
		if result.Status {
			status = "Success"
			success = success + 1
		}

		if result.Status || showFailed {
			if outputFormat == "plain" {
				fmt.Printf("\nURL:\t\t%v\nMethod:\t\t%v\nData:\t\t%v\nResult:\t\t\033[1m%v\033[m\n", result.Request.url, result.Request.method, result.Request.data, status)
			}

		}

		if outputFormat == "json" {
			enc, err := json.Marshal(results)
			if err != nil {
				fmt.Println(err)
			}

			fmt.Println(string(enc))
		}

	}

	if outputFormat == "plain" {
		fmt.Printf("\nTotal success: \033[1m%v\033[m\n", success)
	}

}

func Run(payloads io.Reader, req Request, results chan Result) {

	startPhantomEngine()
	defer phantom.Exit()

	err, vectors := loadPayloads(payloads)
	if err != nil {
		panic(err)
	}

	for _, v := range vectors {
		err, request, result := Evaluate(req, v)

		if err != nil {
			fmt.Printf("Error: %v", err)
		}

		results <- Result{
			Vector:  v,
			Status:  result,
			Request: request,
		}

	}

	close(results)
}

func loadPayloads(handle io.Reader) (error, []XSSVector) {

	var vectors []XSSVector

	scanner := bufio.NewScanner(handle)
	for scanner.Scan() {
		payload := strings.Trim(scanner.Text(), " ")
		if payload == "" {
			continue
		}

		vectors = append(vectors, XSSVector{
			Payload: payload,
		})
	}

	if err := scanner.Err(); err != nil {
		return err, vectors
	}

	return nil, vectors
}
