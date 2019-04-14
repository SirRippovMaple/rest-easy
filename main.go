package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func makeCall(input *Input) {
	client := &http.Client{}
	request, _ := http.NewRequest(input.method, input.url, bytes.NewReader(input.body))
	for _, header := range input.headers {
		index := strings.Index(header, ":")
		request.Header.Add(header[0:index], strings.TrimSpace(header[index+1:]))
	}

	response, _ := client.Do(request)

	fmt.Printf("HTTP %v\n", response.StatusCode)
	for headerKey, headerValue := range response.Header {
		fmt.Printf("%v: %v\n", headerKey, headerValue)
	}

	fmt.Printf("\n")
	responseBody, _ := ioutil.ReadAll(response.Body)
	os.Stdout.Write(responseBody)
}

func main() {
	input := Parse(os.Stdin)
	makeCall(input)
}
