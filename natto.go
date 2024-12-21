package natto

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"golang.org/x/sys/unix"
	"strings"
)

const Version = "0.0.4"

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

func (c *Capsule) SpartanRequest(r io.Reader) (string, string, error) {
	reader := bufio.NewReader(r)
	request, err := reader.ReadString('\n')
	if err != nil {
		return "", "", fmt.Errorf("something went wrong")
	}
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

func (c *Capsule) GeminiRequest(r io.Reader) (string, string, error) {
	reader := bufio.NewReaderSize(r, 1024)
	request, tooLong, err := reader.ReadLine()
	if tooLong {
		return "", "", fmt.Errorf("request is too long")
	}
	if err != nil {
		return "", "", fmt.Errorf("something went wrong")
	}
	req, err := url.Parse(string(request))
	if err != nil {
		return "", "", fmt.Errorf("trouble parsing the url...")
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
	case info.Mode() & 0111 != 0:
		base := filepath.Base(path)
		err := unix.Exec(path, []string{base}, os.Environ())
		if err != nil {
			return c.Panic(50, "something went wrong\n")
		}
	default:
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
	}

	return nil
}
