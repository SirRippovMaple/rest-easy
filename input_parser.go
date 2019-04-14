package main

import (
	"bufio"
	"io"
	"strings"
)

type Input struct {
	method  string
	url     string
	headers []string
	body    []byte
}

type parser struct {
	bodyMode bool
}

func Parse(reader io.Reader) *Input {
	scanner := bufio.NewScanner(reader)
	parser := &parser{bodyMode: false}
	method, url, available := parser.parseFirstLine(scanner)
	if !available {
		panic("No first line")
	}
	headers, available := parser.parseHeaders(scanner)
	body := parser.parseBody(scanner)

	return &Input{*method, *url, headers, body}
}

func (parser *parser) parseFirstLine(scanner *bufio.Scanner) (*string, *string, bool) {
	available := scanner.Scan()
	if !available {
		return nil, nil, false
	}

	line := scanner.Text()
	if space := strings.Index(line, " "); space > 0 {
		method := line[0:space]
		url := line[space+1:]
		return &method, &url, available
	}

	return nil, nil, false
}

func (parser *parser) parseHeaders(scanner *bufio.Scanner) (headers []string, available bool) {
	available = scanner.Scan()
	for available {
		line := scanner.Text()
		if len(line) == 0 {
			break
		}

		headers = append(headers, line)
		available = scanner.Scan()
	}

	return headers, available
}

func (parser *parser) parseBody(scanner *bufio.Scanner) []byte {
	parser.bodyMode = true
	scanner.Scan()
	return scanner.Bytes()
}

func (parser *parser) splitFunc(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if parser.bodyMode {
		return bodySplitFunc(data, atEOF)
	}
	return bufio.ScanLines(data, atEOF)
}

func bodySplitFunc(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	if !atEOF {
		return 0, nil, nil
	}

	return len(data), data, nil
}
