package main

import (
	"bufio"
	"strings"
	"testing"
)

func TestNoReplaceUrl(t *testing.T) {
	url := "https://example.com/"
	variables := make(map[string]string)
	replacedUrl := replaceUrl(url, variables)

	if url != *replacedUrl {
		t.Error("Url does not match")
	}
}

func TestReplaceQueryString(t *testing.T) {
	url := "https://example.com?q={{q}}"
	variables := make(map[string]string)
	variables["q"] = "test"
	replacedUrl := replaceUrl(url, variables)

	if *replacedUrl != "https://example.com?q=test" {
		t.Errorf("Url does not match. Got %v", *replacedUrl)
	}
}

func TestReplaceEmptyQueryString(t *testing.T) {
	url := "https://example.com?q={{q}}"
	variables := make(map[string]string)
	replacedUrl := replaceUrl(url, variables)

	if *replacedUrl != "https://example.com" {
		t.Errorf("Url does not match. Got %v", *replacedUrl)
	}
}

// Scanner tests
func recogniseToken(scanner *bufio.Scanner, expectedToken string, t *testing.T) {
	if !scanner.Scan() {
		t.Errorf("Expected the '%v' token to be available", expectedToken)
	}
	if scanner.Text() != expectedToken {
		t.Errorf("Expected the '%v' token at this point, got '%v' instead.", expectedToken, scanner.Text())
	}
}

func recogniseEof(scanner *bufio.Scanner, t *testing.T) {
	if scanner.Scan() {
		t.Errorf("Expected no more tokens, but got '%v' instead,", scanner.Text())
	}
}

func TestScannerRecognisesOpen(t *testing.T) {
	s := "text{{"
	scanner := bufio.NewScanner(strings.NewReader(s))
	scanner.Split(scan)

	recogniseToken(scanner, "text", t)
	recogniseToken(scanner, "{{", t)
	recogniseEof(scanner, t)
}

func TestScannerRecognisesClose(t *testing.T) {
	s := "text}}"
	scanner := bufio.NewScanner(strings.NewReader(s))
	scanner.Split(scan)

	recogniseToken(scanner, "text", t)
	recogniseToken(scanner, "}}", t)
	recogniseEof(scanner, t)
}

func TestFullVariable(t *testing.T) {
	s := "{{text}}"
	scanner := bufio.NewScanner(strings.NewReader(s))
	scanner.Split(scan)

	recogniseToken(scanner, "{{", t)
	recogniseToken(scanner, "text", t)
	recogniseToken(scanner, "}}", t)
	recogniseEof(scanner, t)
}