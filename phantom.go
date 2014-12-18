package main

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/urturn/go-phantomjs"
)

var phantom *phantomjs.Phantom

var evaluateByVar = `function(done) {
	var page = require('webpage').create();

	page.onError = function() {};

	page.onInitialized = function() {

		page.evaluate(function() {
			var block = document.createElement('script');
			block.type = 'text/javascript';
			block.innerHTML = 'xssm = function xssm() { window.xssmap = true; }';
			var head = document.getElementsByTagName('head')[0];

			head.insertBefore(block, head.firstChild); // Add first
		});

	};

	var url = atob('%s');
	var method = '%s';
	var data = atob('%s');

	page.open(url, method, data, function (status) {
		var xssmap = page.evaluate(function() {
			return (typeof window.xssmap !== 'undefined');
		});
		page.close();		
		setTimeout(function() { done(xssmap); }, 0);
	});
}`

func startPhantomEngine() {
	var err error
	phantom, err = phantomjs.Start()

	if err != nil {
		panic(err)
	}
}

func toBase64(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

func Evaluate(req Request, vector XSSVector) (error, Request, bool) {
	var result interface{}

	// Replace {XSS} placeholder with payload
	if req.method == "GET" {
		req.url = strings.Replace(req.url, "{XSS}", vector.Payload, -1)
	} else {
		req.data = strings.Replace(req.data, "{XSS}", vector.Payload, -1)
	}

	var phantomFunc = fmt.Sprintf(evaluateByVar, toBase64(req.url), req.method, toBase64(req.data))

	err := phantom.Run(phantomFunc, &result)
	if err != nil {
		return err, req, false
	}

	return nil, req, result.(bool)
}
