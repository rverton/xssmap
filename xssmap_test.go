package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDOMXSS(t *testing.T) {

	var payloads = `"><script>xssm()</script>
<img src=x onerror=xssm()>`

	var testpage = `
<html>
<head>
</head>
<body>
<script>
var pos=document.URL.indexOf("xss=")+4;
document.write(decodeURIComponent(document.URL.substring(pos,document.URL.length)));
</script>
</body>
</html>
`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, testpage)
	}))
	defer ts.Close()

	c := make(chan Result)
	var results []Result

	req := Request{
		url:    fmt.Sprintf("%v/#xss={XSS}", ts.URL),
		method: "GET",
	}

	go Run(strings.NewReader(payloads), req, c)

	success := 0

	for res := range c {
		results = append(results, res)
		if res.Status {
			success = success + 1
		}
	}

	if success != 2 {
		t.Error("Tested payloads not identified.")
	}
}
