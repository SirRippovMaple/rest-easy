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
	available := scanner.Scan()
	variables, available := parser.parseFrontMatter(scanner, available)
	method, url, available := parser.parseFirstLine(scanner, variables, available)
	if !available {
		panic("No first line")
	}
	headers, available := parser.parseHeaders(scanner, variables)
	body := parser.parseBody(scanner)

	return &Input{*method, *url, headers, body}
}

func (parser *parser) parseFrontMatter(scanner *bufio.Scanner, available bool) (map[string]string, bool) {
	variables := make(map[string]string)

	if scanner.Text() != "---" {
		return make(map[string]string), available
	}

	for scanner.Scan() && scanner.Text() != "---" {
		line := scanner.Text()
		colonIndex := strings.Index(line, ":")
		key := strings.TrimSpace(line[:colonIndex])
		value := strings.TrimSpace(line[colonIndex+1:])
		variables[key] = value
	}

	return variables, scanner.Scan()
}

func (parser *parser) parseFirstLine(scanner *bufio.Scanner, variables map[string]string, available bool) (method *string, url *string, returnAvailable bool) {
	if !available {
		return nil, nil, false
	}

	line := scanner.Text()
	if space := strings.Index(line, " "); space > 0 {
		method := line[0:space]
		url := replaceVariables(line[space+1:], variables)
		return &method, &url, available
	}

	return nil, nil, false
}

func (parser *parser) parseHeaders(scanner *bufio.Scanner, variables map[string]string) (headers []string, available bool) {
	available = scanner.Scan()
	for available {
		line := scanner.Text()
		if len(line) == 0 {
			break
		}

		headers = append(headers, replaceVariables(line, variables))
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

// Probably very slow, but this works for now
func replaceVariables(input string, variables map[string]string) string {
	current := input
	for k, v := range variables {
		token := "{{" + k + "}}"
		current = strings.Replace(current, token, v, -1)
	}

	return current
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
