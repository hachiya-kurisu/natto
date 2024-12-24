package natto

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const Version = "0.0.7"

const (
	Input               = 10
	SensitiveInput      = 11
	Success             = 20
	TemporaryRedirect   = 30
	PermanentRedirect   = 31
	TemporaryFailure    = 40
	ServerUnavailable   = 41
	CGIError            = 42
	ProxyError          = 43
	SlowDown            = 44
	PermanentFailure    = 50
	NotFound            = 51
	Gone                = 52
	ProxyRequestRefused = 53
	BadRequest          = 59
)

const (
	Redirect    = 3
	ClientError = 4
	ServerError = 5
)

type Protocol int

const (
	Gemini Protocol = iota
	Spartan
)

type Failure struct {
	status  int
	message string
}

func (f *Failure) Error() string {
	return f.message
}

var Types = map[string]string{
	".gmi":  "text/gemini",
	".txt":  "text/plain",
	".jpg":  "image/jpeg",
	".png":  "image/png",
	".jxl":  "image/jxl",
	".webp": "image/webp",
	".mp3":  "audio/mpeg",
	".m4a":  "audio/mp4",
	".mp4":  "video/mp4",
}

type CgiHandler func(string, []string, []string) error

type Capsule struct {
	Path       string
	Writer     io.Writer
	Protocol   Protocol
	CgiHandler CgiHandler
}

func (c *Capsule) ProtocolName() string {
	switch c.Protocol {
	case Spartan:
		return "spartan"
	default:
		return "gemini"
	}
}

func (c *Capsule) SpartanRequest(request string) (string, string, error) {
	request = strings.TrimSpace(request)
	components := strings.SplitN(request, " ", 3)
	if len(components) != 3 {
		return c.Fail(ClientError, "bad request")
	}
	host, path, contentLength := components[0], components[1], components[2]

	if path[0] != '/' {
		return c.Fail(ClientError, "bad request - missing /")
	}

	length, err := strconv.Atoi(contentLength)
	if err != nil {
		return c.Fail(ClientError, "bad request - invalid content length")
	}
	if length > 0 {
		return c.Fail(ServerError, "no data block support yet")
	}

	return host, path, nil
}

func (c *Capsule) Fail(status int, message string) (string, string, error) {
	return "", "", &Failure{status, message}
}

func (c *Capsule) GeminiRequest(request string) (string, string, error) {
	if len(request) > 1024 {
		return c.Fail(BadRequest, "bad request: too long")
	}
	req, err := url.Parse(strings.TrimSpace(request))
	if err != nil {
		return c.Fail(BadRequest, "invalid url")
	}
	err = c.Validate(req)
	if err != nil {
		return c.Fail(BadRequest, err.Error())
	}
	return req.Host, req.Path, nil
}

func (c *Capsule) Validate(req *url.URL) error {
	if !req.IsAbs() {
		return fmt.Errorf("not an absolute url")
	}
	if req.Fragment != "" {
		return fmt.Errorf("fragments not allowed")
	}
	if req.User != nil {
		return fmt.Errorf("userinfo not allowed")
	}
	if req.Scheme != "gemini" {
		return fmt.Errorf("this is a gemini server")
	}
	return nil
}

func (c *Capsule) Panic(failure *Failure) error {
	status := failure.status
	if c.Protocol == Spartan && status > 10 {
		status /= 10
	}
	fmt.Fprintf(c.Writer, "%d %s\r\n", status, failure.Error())
	return failure
}

func (c *Capsule) Header(status int, info string) {
	switch c.Protocol {
	case Spartan:
		fmt.Fprintf(c.Writer, "%d %s\r\n", status/10, info)
	default:
		fmt.Fprintf(c.Writer, "%d %s\r\n", status, info)
	}
}

func (c *Capsule) Cgi(path string, handler CgiHandler) error {
	os.Setenv("GATEWAY_INTERFACE", "CGI/1.1")
	os.Setenv("SERVER_PROTOCOL", c.ProtocolName())
	base := filepath.Base(path)
	err := handler(path, []string{base}, os.Environ())
	if err != nil {
		return c.Panic(&Failure{CGIError, "something went wrong"})
	}
	return nil
}

func (c *Capsule) Serve(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return c.Panic(&Failure{NotFound, "not found"})
	}
	mime := Types[filepath.Ext(path)]
	if mime == "" {
		mime = "application/octet-stream"
	}
	c.Header(Success, mime)
	io.Copy(c.Writer, f)
	return nil
}

func (c *Capsule) Request(host, path string) error {
	if c.Writer == nil {
		c.Writer = ioutil.Discard
	}
	if path == "" {
		path = "/"
	}
	if path[len(path)-1] == '/' {
		path = path + "index.gmi"
	}

	path = "." + path
	info, err := os.Stat(path)
	if err != nil {
		return c.Panic(&Failure{NotFound, "not found"})
	}

	switch {
	case info.Mode()&0111 != 0 && c.CgiHandler != nil:
		return c.Cgi(path, c.CgiHandler)
	default:
		return c.Serve(path)
	}
}
