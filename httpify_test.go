package httpify_test

import (
	"net/url"
	"testing"

	"github.com/Melidee/httpify"
)

func TestReadRequest(t *testing.T) {
	req, err := httpify.ReadRequest([]byte("GET / HTTP/1.1\r\nHost: localhost:8080\r\n\r\n"))
	if err != nil {
		t.Error(err)
	}
	if req.Method() != "GET" {
		t.Error("invalid method")
	}
	if req.Url().String() != "/" {
		t.Error("invalid resource")
	}
	if req.ProtoMajor() != 1 || req.ProtoMinor() != 1 {
		t.Error("invalid protocol version")
	}
	if req.Headers()["Host"] != "localhost:8080" {
		t.Error("invalid header")
	}
	if req.Body() != "" {
		t.Error("invalid body")
	}
}

func TestRequestString(t *testing.T) {
	url, _ := url.Parse("/")
	req := httpify.NewRequest("GET", url, 1, 1, map[string]string{"Host": "localhost:8080"}, "Hello World!")
	if req.String() != "GET / HTTP/1.1\r\nHost: localhost:8080\r\n\r\nHello World!" {
		t.Error("invalid string")
	}
}

func TestReadResponse(t *testing.T) {
	res, err := httpify.ReadResponse([]byte("HTTP/1.1 200 OK\r\nContent-Length: 0\r\n\r\n"))
	if err != nil {
		t.Error(err)
	}
	if res.ProtoMajor() != 1 || res.ProtoMinor() != 1 {
		t.Error("invalid protocol version")
	}
	if res.StatusCode() != 200 {
		t.Error("invalid status code")
	}
	if res.Status() != "200 OK" {
		t.Error("invalid status")
	}
	if res.Headers()["Content-Length"] != "0" {
		t.Error("invalid header")
	}
	if res.Body() != "" {
		t.Error("invalid body")
	}
}

func TestResponseString(t *testing.T) {
	res := httpify.NewResponse(1, 1, 200, map[string]string{"Content-Length": "0"}, "Hello World!")
	if res.String() != "HTTP/1.1 200 OK\r\nContent-Length: 0\r\n\r\nHello World!" {
		t.Error("invalid string")
	}
}