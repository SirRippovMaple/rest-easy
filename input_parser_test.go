package main

import (
	"bytes"
	"os"
	"testing"
)

func assertInput(t *testing.T, expected *Input, actual *Input) {
	if actual.method != expected.method {
		t.Errorf("Unexpected method. Expected %v, but got %v", expected.method, actual.method)
	}

	if actual.url != expected.url {
		t.Errorf("Received the incorrect url. Expected %v, but got %v", expected.url, actual.url)
	}

	if len(actual.headers) != len(expected.headers) {
		t.Errorf("Received the incorrect number of headers. Expected %v, but got %v", len(expected.headers), len(actual.headers))
	}

	for index, actualHeader := range actual.headers {
		expectedHeader := expected.headers[index]

		if actualHeader != expectedHeader {
			t.Errorf("Header %v: expected %v, but got %v", index, expectedHeader, actualHeader)
		}
	}

	if bytes.Compare(expected.body, actual.body) != 0 {
		t.Errorf("Received incorrect body, expected %v, but got %v", expected.body, actual.body)
	}
}

func TestSimpleGetRequest(t *testing.T) {
	expectedInput := &Input{
		method:  "GET",
		url:     "https://example.com/",
		headers: []string{"Authorization: Bearer FAKE_TOKEN"},
	}

	reader, _ := os.Open("TestFiles/simple-get.http")
	input := Parse(reader, make(map[string]string))

	assertInput(t, expectedInput, input)
}

func TestSimplePostRequest(t *testing.T) {
	expectedInput := &Input{
		method: "POST",
		url:    "https://example.com/",
		headers: []string{
			"Authorization: Bearer FAKE_TOKEN",
			"Content-Type: text/plain",
		},
		body: []byte("This is the body"),
	}

	reader, _ := os.Open("TestFiles/simple-post.http")
	input := Parse(reader, make(map[string]string))

	assertInput(t, expectedInput, input)
}

func TestTemplateGetRequest(t *testing.T) {
	expectedInput := &Input {
		method: "GET",
		url: "https://example.com/1",
		headers: []string{
			"Authorization: Bearer FAKE_TOKEN_1",
		},
	}

	reader, _ := os.Open("TestFiles/template-get.http")
	input := Parse(reader, make(map[string]string))

	assertInput(t, expectedInput, input)
}