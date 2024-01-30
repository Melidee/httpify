package httpify_test

import (
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
	if req.Resource().String() != "/" {
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