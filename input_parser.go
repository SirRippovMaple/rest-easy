package main

import (
	"bufio"
	"io"
	"strings"
)

type Input struct {
	method  string
	originalUrl string
	url     string
	headers []string
	body    []byte
}

type parser struct {
	bodyMode bool
}

func Parse(reader io.Reader, variables map[string]string) *Input {
	scanner := bufio.NewScanner(reader)
	parser := &parser{bodyMode: false}
	available := scanner.Scan()
	available = parser.parseFrontMatter(scanner, available, variables)
	method, url, originalUrl, available := parser.parseFirstLine(scanner, variables, available)
	if !available {
		panic("No first line")
	}
	headers, available := parser.parseHeaders(scanner, variables)
	body := parser.parseBody(scanner)

	return &Input{*method, *originalUrl, *url, headers, body}
}

func (parser *parser) parseFrontMatter(scanner *bufio.Scanner, available bool, variables map[string]string) bool {
	if scanner.Text() != "---" {
		return available
	}

	for scanner.Scan() && scanner.Text() != "---" {
		line := scanner.Text()
		colonIndex := strings.Index(line, ":")
		key := strings.TrimSpace(line[:colonIndex])
		value := strings.TrimSpace(line[colonIndex+1:])
		if _, exists := variables[key]; !exists {
			variables[key] = value
		}
	}

	return scanner.Scan()
}

func (parser *parser) parseFirstLine(scanner *bufio.Scanner, variables map[string]string, available bool) (method , url, originalUrl *string, returnAvailable bool) {
	if !available {
		return nil, nil, nil, false
	}

	line := scanner.Text()
	if space := strings.Index(line, " "); space > 0 {
		method := line[0:space]
		originalUrl := line[space+1:]
		url := replaceUrl(originalUrl, variables)
		return &method, url, &originalUrl, available
	}

	return nil, nil, nil,false
}

func (parser *parser) parseHeaders(scanner *bufio.Scanner, variables map[string]string) (headers []string, available bool) {
	available = scanner.Scan()
	for available {
		line := scanner.Text()
		if len(line) == 0 {
			break
		}

		headers = append(headers, *replaceVariables(line, variables))
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
