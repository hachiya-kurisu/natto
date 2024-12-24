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

const Version = "0.0.5"

type Protocol int

const (
	Gemini Protocol = iota
	Spartan
)

var Types = map[string]string{
	".gmi": "text/gemini",
	".txt": "text/plain",
	".jpg": "image/jpeg",
	".m4a": "audio/mp4",
	".mp4": "video/mp4",
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
		return "", "", fmt.Errorf("invalid request: need host, path and length")
	}
	host, path, contentLength := components[0], components[1], components[2]

	if path[0] != '/' {
		return "", "", fmt.Errorf("invalid request: path has to begin with a /")
	}

	length, err := strconv.Atoi(contentLength)
	if err != nil {
		return "", "", fmt.Errorf("invalid content length")
	}
	if length > 0 {
		return "", "", fmt.Errorf("no data block support (yet)")
	}

	return host, path, nil
}

func (c *Capsule) GeminiRequest(request string) (string, string, error) {
	if len(request) > 1024 {
		return "", "", fmt.Errorf("request is too long")
	}
	req, err := url.Parse(strings.TrimSpace(request))
	if err != nil {
		return "", "", fmt.Errorf("trouble parsing the url")
	}
	err = c.Validate(req)
	if err != nil {
		return "", "", err
	}
	return req.Host, req.Path, nil
}

func (c *Capsule) Validate(req *url.URL) error {
	if !req.IsAbs() {
		return fmt.Errorf("invalid request: not an absolute url")
	}
	if req.Fragment != "" {
		return fmt.Errorf("invalid request: fragments not allowed")
	}
	if req.User != nil {
		return fmt.Errorf("invalid request: userinfo not allowed")
	}
	if req.Scheme != "gemini" {
		return fmt.Errorf("invalid request: this is a gemini server")
	}
	return nil
}

func (c *Capsule) Panic(status int, response string) error {
	switch c.Protocol {
	case Spartan:
		fmt.Fprintf(c.Writer, "%d\r\n", status/10)
	case Gemini:
		fmt.Fprintf(c.Writer, "%d\r\n", status)
	}
	fmt.Fprintln(c.Writer, response)
	return fmt.Errorf(response)
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
	// unix.Exec(path, []string{base}, os.Environ())
	if err != nil {
		return c.Panic(50, "something went wrong\n")
	}
	return nil
}

func (c *Capsule) Serve(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return c.Panic(40, "not found\n")
	}
	mime := Types[filepath.Ext(path)]
	if mime == "" {
		mime = "application/octet-stream"
	}
	c.Header(20, mime)
	io.Copy(c.Writer, f) // ignore errors until we have proper logging
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
		return c.Panic(40, "not found\n")
	}

	switch {
	case info.Mode()&0111 != 0 && c.CgiHandler != nil:
		return c.Cgi(path, c.CgiHandler)
	default:
		return c.Serve(path)
	}
}
