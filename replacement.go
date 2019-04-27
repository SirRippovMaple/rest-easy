package main

import (
	"bufio"
	"bytes"
	"net/url"
	"strings"
)

func replaceUrl(originalUrl string, variables map[string]string) *string {
	u, _ := url.Parse(originalUrl)

	newHost := replaceVariables(u.Host, variables)
	var newQs = make(url.Values)
	for qk, qv := range u.Query() {
		for _, sqv := range qv {
			newQv := replaceVariables(sqv, variables)
			if newQv != nil {
				newQs.Add(qk, *newQv)
			}
		}
	}

	newUrl := &url.URL{Host: *newHost, Scheme: u.Scheme, Path:u.Path, RawQuery: newQs.Encode()}
	newUrlString := newUrl.String()
	return &newUrlString
}

// This works as long as the string is valid. We probably need a state machine to recognise invalid strings
func replaceVariables(input string, variables map[string]string) *string {
	var output string
	reader := strings.NewReader(input)
	scanner := bufio.NewScanner(reader)
	scanner.Split(scan)
	inVar := false

	for scanner.Scan() {
		text := scanner.Text()

		switch text {
		case "{{":
			inVar = true
			break
		case "}}":
			if !inVar {
				output += text
			}
			inVar = false
			break
		default:
			if inVar {
				replacement := variables[text]
				output += replacement
			} else {
				output += text
			}
		}
	}

	if len(output) == 0 {
		return nil
	}
	return &output
}

func scan(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	if i := bytes.IndexByte(data, '{'); i > 0 {
		if len(data) > i && data[i+1] == '{' {
			return i, data[:i], nil
		}
	}

	if len(data) >= 2 && data[0] == '{' && data[1] == '{' {
		return 2, data[0:2], nil
	}

	if i := bytes.IndexByte(data, '}'); i > 0 {
		if len(data) > i && data[i+1] == '}' {
			return i, data[:i], nil
		}
	}

	if len(data) >= 2 && data[0] == '}' && data[1] == '}' {
		return 2, data[0:2], nil
	}
	if !atEOF {
		return 0, nil, nil
	}

	return len(data), data, nil
}