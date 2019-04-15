package main

import (
	"bytes"
	"flag"
	"fmt"
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

func main() {
	input := flag.String("input", "", "The request filer")
	output := flag.String("output", "", "The file to write the response output")
	flag.Parse()
	
	var requestReader io.Reader
	var responseWriter io.Writer

	if len(*input) > 0 {
		requestReader, _ = os.Open(*input)
	} else {
		stats, _ := os.Stdin.Stat()
		if stats.Size() > 0 {
			requestReader = os.Stdin
		} else {
			flag.Usage()
			return
		}
	}

	if len(*output) > 0 {
		responseWriter = os.Stdout
	} else {
		responseWriter = os.Stdout
	}

	request := Parse(requestReader)
	makeCall(request, responseWriter)
}
