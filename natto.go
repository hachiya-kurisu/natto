package natto

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
)

const Version = "0.0.4"

type Result int

const (
	Ok Result = iota
	Oops
)

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

type Capsule struct {
	Path     string
	Writer   io.Writer
	Protocol Protocol
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

func (c *Capsule) Panic(status int, response string) {
	switch c.Protocol {
	case Spartan:
		fmt.Fprintf(c.Writer, "%d\r\n", status/10)
	case Gemini:
		fmt.Fprintf(c.Writer, "%d\r\n", status)
	}
	fmt.Fprintln(c.Writer, response)
}

func (c *Capsule) Header(status int, info string) {
	switch c.Protocol {
	case Spartan:
		fmt.Fprintf(c.Writer, "%d %s\r\n", status/10, info)
	default:
		fmt.Fprintf(c.Writer, "%d %s\r\n", status, info)
	}
}

func (c *Capsule) Request(host, path string) Result {
	if c.Writer == nil {
		c.Writer = ioutil.Discard
	}
	if path == "" {
		path = "/"
	}
	if path[len(path)-1] == '/' {
		path = path + "index.gmi"
	}

	f, err := os.Open("." + path)
	if err != nil {
		c.Panic(40, "not found\n")
		return Oops
	}

	mime := Types[filepath.Ext(path)]
	if mime == "" {
		mime = "application/octet-stream"
	}

	c.Header(20, mime)
	io.Copy(c.Writer, f) // ignore errors until we have proper logging

	return Ok
}
