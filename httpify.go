package httpify

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type Request struct {
	method     string
	url        *url.URL
	protoMajor int
	protoMinor int
	headers    map[string]string
	body       string
}

func NewRequest(method string, url *url.URL, protoMajor int, protoMinor int, headers map[string]string, body string) *Request {
	return &Request{
		method:     method,
		url:        url,
		protoMajor: protoMajor,
		protoMinor: protoMinor,
		headers:    headers,
		body:       body,
	}
}

// TODO: more in depth errors
func ReadRequest(data []byte) (req *Request, err error) {
	if len(data) == 0 {
		return nil, errors.New("no data provided")
	}
	req_text := string(data)
	lines := strings.Split(req_text, "\r\n")
	first_line := strings.SplitN(lines[0], " ", 3)
	method, raw_url, proto := first_line[0], first_line[1], first_line[2]
	_, ok := HTTP_METHOD[method]
	if !ok {
		return nil, errors.New("invalid HTTP method")
	}
	url, err := url.Parse(raw_url)
	if err != nil {
		return nil, errors.New("invalid url")
	}
	var protoMajor, protoMinor int
	if proto == "HTTP/1.1" {
		protoMajor, protoMinor = 1, 1
	} else {
		return nil, errors.New("only HTTP/1.1 parsing is supported")
	}

	headers := map[string]string{}
	var l int
	for l = 1; lines[l] != ""; l++ {
		header_parts := strings.SplitN(lines[l], ": ", 2)
		headers[header_parts[0]] = header_parts[1]
	}
	body := strings.Join(lines[l+1:], "\r\n")

	return &Request{
		method:     method,
		url:        url,
		protoMajor: protoMajor,
		protoMinor: protoMinor,
		headers:    headers,
		body:       body,
	}, nil
}

func (req *Request) String() string {
	return req.Method() + " " + req.Url().String() + " " + req.Protocol() + "\r\n" + req.HeadersString() + "\r\n" + req.Body()
}

func (req *Request) Dump() []byte {
	return []byte(req.String())
}

func (req *Request) Method() string {
	return req.method
}

func (req *Request) SetMethod(method string) {
	req.method = method
}

func (req *Request) Url() *url.URL {
	return req.url
}

func (req *Request) SetUrl(url *url.URL) {
	req.url = url
}

func (req *Request) Resource() *url.URL {
	return req.url
}

func (req *Request) SetResource(resource string) {
	req.url, _ = url.Parse(resource)
}

func (req *Request) ProtoMajor() int {
	return req.protoMajor
}

func (req *Request) SetProtoMajor(protoMajor int) {
	req.protoMajor = protoMajor
}

func (req *Request) ProtoMinor() int {
	return req.protoMinor
}

func (req *Request) SetProtoMinor(protoMinor int) {
	req.protoMinor = protoMinor
}

func (req *Request) Protocol() string {
	return "HTTP/" + fmt.Sprint(req.protoMajor) + "." + fmt.Sprint(req.protoMinor)
}

func (req *Request) SetProtocol(protoMajor int, protoMinor int) {
	req.protoMajor = protoMajor
	req.protoMinor = protoMinor
}

func (req *Request) Headers() map[string]string {
	return req.headers
}

func (req *Request) SetHeader(key string, value string) {
	req.headers[key] = value
}

func (req *Request) HeadersString() string {
	headers := ""
	for key, value := range req.headers {
		headers += key + ": " + value + "\r\n"
	}
	return headers
}

func (req *Request) Body() string {
	return req.body
}

func (req *Request) SetBody(body string) {
	req.body = body
}

type Response struct {
	protoMajor int
	protoMinor int
	statusCode int
	headers    map[string]string
	body       string
}

func NewResponse(protoMajor int, protoMinor int, statusCode int, headers map[string]string, body string) *Response {
	return &Response{
		protoMajor: protoMajor,
		protoMinor: protoMinor,
		statusCode: statusCode,
		headers:    headers,
		body:       body,
	}
}

func ReadResponse(data []byte) (res *Response, err error) {
	if len(data) == 0 {
		return nil, errors.New("no data provided")
	}
	res_text := string(data)
	lines := strings.Split(res_text, "\r\n")
	first_line := strings.SplitN(lines[0], " ", 2)
	var protoMajor, protoMinor int
	if first_line[0] == "HTTP/1.1" {
		protoMajor, protoMinor = 1, 1
	} else {
		return nil, errors.New("only HTTP/1.1 parsing is supported")
	}
	statusCode, err := strconv.Atoi(first_line[1][0:3])
	if err != nil || HTTP_STATUS[statusCode] == "" {
		return nil, errors.New("invalid status code")
	}

	headers := map[string]string{}
	var l int
	for l = 1; lines[l] != ""; l++ {
		header_parts := strings.SplitN(lines[l], ": ", 2)
		headers[header_parts[0]] = header_parts[1]
	}
	body := strings.Join(lines[l+1:], "\r\n")

	return &Response{
		protoMajor: protoMajor,
		protoMinor: protoMinor,
		statusCode: statusCode,
		headers:    headers,
		body:       body,
	}, nil
}

func (req *Response) String() string {
	return req.Protocol() + " " + req.Status() + "\r\n" + req.HeadersString() + "\r\n" + req.Body()
}

func (req *Response) ProtoMajor() int {
	return req.protoMajor
}

func (req *Response) SetProtoMajor(protoMajor int) {
	req.protoMajor = protoMajor
}

func (req *Response) ProtoMinor() int {
	return req.protoMinor
}

func (req *Response) SetProtoMinor(protoMinor int) {
	req.protoMinor = protoMinor
}

func (req *Response) Protocol() string {
	return "HTTP/" + fmt.Sprint(req.protoMajor) + "." + fmt.Sprint(req.protoMinor)
}

func (req *Response) SetProtocol(protoMajor int, protoMinor int) {
	req.protoMajor = protoMajor
	req.protoMinor = protoMinor
}

func (req *Response) StatusCode() int {
	return req.statusCode
}

func (req *Response) SetStatusCode(statusCode int) {
	req.statusCode = statusCode
}

func (req *Response) Status() string {
	return HTTP_STATUS[req.statusCode]
}

func (req *Response) SetStatus(statusCode int) {
	req.statusCode = statusCode
}

func (req *Response) Headers() map[string]string {
	return req.headers
}

func (req *Response) SetHeader(key string, value string) {
	req.headers[key] = value
}

func (req *Response) HeadersString() string {
	headers := ""
	for key, value := range req.headers {
		headers += key + ": " + value + "\r\n"
	}
	return headers
}

func (req *Response) Body() string {
	return req.body
}

func (req *Response) SetBody(body string) {
	req.body = body
}

var HTTP_METHOD = map[string]string{
	"GET":     "GET",
	"HEAD":    "HEAD",
	"POST":    "POST",
	"PUT":     "PUT",
	"DELETE":  "DELETE",
	"CONNECT": "CONNECT",
	"OPTIONS": "OPTIONS",
	"TRACE":   "TRACE",
	"PATCH":   "PATCH",
}

var HTTP_STATUS = map[int]string{
	100: "100 Continue",
	101: "101 Switching Protocols",
	102: "102 Processing",
	103: "103 Early Hints",
	200: "200 OK",
	201: "201 Created",
	202: "202 Accepted",
	203: "203 Non-Authoritative Information",
	204: "204 No Content",
	205: "205 Reset Content",
	206: "206 Partial Content",
	207: "207 Multi-Status",
	208: "208 Already Reported",
	226: "226 IM Used",
	300: "300 Multiple Choices",
	301: "301 Moved Permanently",
	302: "302 Found",
	303: "303 See Other",
	304: "304 Not Modified",
	307: "307 Temporary Redirect",
	308: "308 Permanent Redirect",
	400: "400 Bad Request",
	401: "401 Unauthorized",
	402: "402 Payment Required",
	403: "403 Forbidden",
	404: "404 Not Found",
	405: "405 Method Not Allowed",
	406: "406 Not Acceptable",
	407: "407 Proxy Authentication Required",
	408: "408 Request Timeout",
	409: "409 Conflict",
	410: "410 Gone",
	411: "411 Length Required",
	412: "412 Precondition Failed",
	413: "413 Payload Too Large",
	414: "414 URI Too Long",
	415: "415 Unsupported Media Type",
	416: "416 Range Not Satisfiable",
	417: "417 Expectation Failed",
	418: "418 I'm a teapot",
	421: "421 Misdirected Request",
	422: "422 Unprocessable Content",
	423: "423 Locked",
	424: "424 Failed Dependency",
	425: "425 Too Early",
	426: "426 Upgrade Required",
	428: "428 Precondition Required",
	429: "429 Too Many Requests",
	431: "431 Request Header Fields Too Large",
	451: "451 Unavailable For Legal Reasons",
	500: "500 Internal Server Error",
	501: "501 Not Implemented",
	502: "502 Bad Gateway",
	503: "503 Service Unavailable",
	504: "504 Gateway Timeout",
	505: "505 HTTP Version Not Supported",
	506: "506 Variant Also Negotiates",
	507: "507 Insufficient Storage",
	508: "508 Loop Detected",
	510: "510 Not Extended",
	511: "511 Network Authentication Required",
}
