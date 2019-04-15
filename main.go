package main

import (
	"bytes"
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"io"
	"net/http"
	"os"
	"strings"
)

func makeCall(input *Input, writer io.Writer) {
	client := &http.Client{}
	request, _ := http.NewRequest(input.method, input.url, bytes.NewReader(input.body))
	for _, header := range input.headers {
		index := strings.Index(header, ":")
		request.Header.Add(header[0:index], strings.TrimSpace(header[index+1:]))
	}

	response, _ := client.Do(request)

	fmt.Fprintf(writer, "HTTP %v\n", response.StatusCode)
	for headerKey, headerValue := range response.Header {
		fmt.Fprintf(writer, "%v: %v\n", headerKey, headerValue)
	}

	fmt.Fprintf(writer, "\n")
	io.Copy(writer, response.Body)
}

var (
	app = kingpin.New("rest-easy", "A command-line http execution utility.")
	run = app.Command("run", "Run a request")
	runInput = run.Arg("input", "Input file. This can be omitted if piping from stdin.").String()
	runOutput = run.Flag("output", "Output file. If this is omitted, then the response is written to stdout.").Short('o').String()
)

func main() {
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case run.FullCommand():
		executeRun(runInput, runOutput)
	}
}

func executeRun(runInput, runOutput *string) {
	var requestReader io.Reader
	var responseWriter io.Writer

	if len(*runInput) > 0 {
		requestReader, _ = os.Open(*runInput)
	} else {
		stats, _ := os.Stdin.Stat()
		if stats.Size() > 0 {
			requestReader = os.Stdin
		} else {
			app.Usage([]string {"run"})
			return
		}
	}

	if len(*runOutput) > 0 {
		responseWriter = os.Stdout
	} else {
		responseWriter = os.Stdout
	}

	request := Parse(requestReader)
	makeCall(request, responseWriter)
}
