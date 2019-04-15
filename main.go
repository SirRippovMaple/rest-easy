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
	run = app.Command("run", "Run a request").Alias("r")
	runInput = run.Arg("input", "Input file. This can be omitted if piping from stdin.").String()
	runOutput = run.Flag("output", "Output file. If this is omitted, then the response is written to stdout.").Short('o').String()
	info = app.Command("info", "Get information about a request").Alias("i")
	infoInput = info.Arg("input", "Input file. This can be omitted if piping from stdin").String()
)

func main() {
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case run.FullCommand():
		executeRun(runInput, runOutput)
	case info.FullCommand():
		executeInfo(infoInput)
	}
}

func executeInfo(infoInput *string) {
	var requestReader io.Reader

	if len(*infoInput) > 0 {
		fileReader, _ := os.Open(*infoInput)
		defer fileReader.Close()
		requestReader = fileReader
	} else {
		stats, _ := os.Stdin.Stat()
		if stats.Size() > 0 {
			requestReader = os.Stdin
		} else {
			app.Usage([]string {"info"})
			return
		}
	}

	variables := make(map[string]string)
	request := Parse(requestReader, variables)

	fmt.Printf("%v %v\n\n", request.method, request.originalUrl)

	fmt.Printf("Variables:\n")
	if len(variables) == 0 {
		fmt.Printf("No variables.\n")
	} else {
		for k, v := range variables {
			fmt.Printf("%v = %v\n", k, v)
		}
	}
}

func executeRun(runInput, runOutput *string) {
	var requestReader io.Reader
	var responseWriter io.Writer

	if len(*runInput) > 0 {
		fileReader, _ := os.Open(*runInput)
		defer fileReader.Close()
		requestReader = fileReader
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

	variables := make(map[string]string)
	request := Parse(requestReader, variables)
	makeCall(request, responseWriter)
}
