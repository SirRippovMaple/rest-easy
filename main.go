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

func makeCall(input *Input, writer io.Writer) error{
	client := &http.Client{}
	request, err := http.NewRequest(input.method, input.url, bytes.NewReader(input.body))
	if err != nil {
		return err
	}
	for _, header := range input.headers {
		index := strings.Index(header, ":")
		request.Header.Add(header[0:index], strings.TrimSpace(header[index+1:]))
	}

	response, err := client.Do(request)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(writer, "HTTP %v\n", response.StatusCode)
	if err != nil {
		return err
	}
	for headerKey, headerValue := range response.Header {
		_, err = fmt.Fprintf(writer, "%v: %v\n", headerKey, headerValue)
		if err != nil {
			return err
		}
	}

	_, err = fmt.Fprintf(writer, "\n")
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, response.Body)
	if err != nil {
		return err
	}

	return nil
}

var (
	app = kingpin.New("rest-easy", "A command-line http execution utility.")
	run = app.Command("run", "Run a request").Alias("r")
	runInput = run.Arg("input", "Input file. This can be omitted if piping from stdin.").String()
	runOutput = run.Flag("output", "Output file. If this is omitted, then the response is written to stdout.").Short('o').String()
	runVariables = run.Flag("set-property", "Sets or overrides a property/variable.").Short('p').StringMap()
	info = app.Command("info", "Get information about a request").Alias("i")
	infoInput = info.Arg("input", "Input file. This can be omitted if piping from stdin").String()
)

func main() {
	var err error

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case run.FullCommand():
		err = executeRun(runInput, runOutput)
		break

	case info.FullCommand():
		err = executeInfo(infoInput)
		break
	}

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, err.Error())
		_, _ = fmt.Fprintf(os.Stderr, "\n")
	}

}

func executeInfo(infoInput *string) error {
	var requestReader io.Reader

	if len(*infoInput) > 0 {
		fileReader, err := os.Open(*infoInput)

		if fileReader != nil {
			defer func() {
				_ = fileReader.Close()
			}()
		}
		if err != nil {
			return err
		}

		requestReader = fileReader
	} else {
		stats, err := os.Stdin.Stat()
		if err != nil {
			return err
		}
		if stats.Size() > 0 {
			requestReader = os.Stdin
		} else {
			app.Usage([]string {"info"})
			return nil
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

	return nil
}

func executeRun(runInput, runOutput *string) error {
	var requestReader io.Reader
	var responseWriter io.Writer

	if len(*runInput) > 0 {
		fileReader, err := os.Open(*runInput)
		if fileReader != nil {
			defer func() {
				_ = fileReader.Close()
			}()
		}
		if err != nil {
			return err
		}

		requestReader = fileReader
	} else {
		stats, _ := os.Stdin.Stat()
		if stats.Size() > 0 {
			requestReader = os.Stdin
		} else {
			app.Usage([]string {"run"})
			return nil
		}
	}

	if len(*runOutput) > 0 {
		responseWriter = os.Stdout
	} else {
		responseWriter = os.Stdout
	}

	request := Parse(requestReader, *runVariables)
	err := makeCall(request, responseWriter)

	return err
}
