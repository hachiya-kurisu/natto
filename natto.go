package natto

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
)

const Version = "0.0.2"

const (
	Ok int = iota + 1
	Oops
)

var Types = map[string]string{
	".gmi": "text/gemini",
	".txt": "text/plain",
	".jpg": "image/jpeg",
	".m4a": "audio/mp4",
	".mp4": "video/mp4",
}

type Capsule struct {
	Path   string
	Writer io.Writer
}

func (c *Capsule) Validate(req *url.URL) {
	if !req.IsAbs() {
		c.Panic(59, "invalid request: not an absolute url")
		os.Exit(1)
	}
	if req.Fragment != "" {
		c.Panic(59, "invalid request: fragments not allowed")
		os.Exit(1)
	}
	if req.User != nil {
		c.Panic(59, "invalid request: userinfo not allowed")
		os.Exit(1)
	}
	if req.Scheme != "gemini" {
		c.Panic(59, "invalid request: this is a gemini server")
		os.Exit(1)
	}
}

func (c *Capsule) Panic(status int, response string) {
	fmt.Fprintf(c.Writer, "%d\r\n", status)
	fmt.Fprintln(c.Writer, response)
}

func (c *Capsule) Request(req *url.URL) int {
	if c.Writer == nil {
		c.Writer = os.Stdout
	}
	if req.Path == "" {
		req.Path = "/"
	}
	if req.Path[len(req.Path)-1] == '/' {
		req.Path = req.Path + "index.gmi"
	}

	f, err := os.Open("." + req.Path)
	if err != nil {
		c.Panic(40, "not found\n")
		return Oops
	}

	mime := Types[filepath.Ext(req.Path)]
	if mime == "" {
		mime = "application/octet-stream"
	}

	fmt.Fprintf(c.Writer, "20 %s\r\n", mime)
	io.Copy(c.Writer, f) // ignore errors until we have proper logging

	return Ok
}
