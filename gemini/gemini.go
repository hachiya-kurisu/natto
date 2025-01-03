package gemini

import (
	"blekksprut.net/natto"
	"fmt"
	"io"
	"io/fs"
	"net/url"
	"os"
	"strings"
)

type Capsule struct {
	Root string
	FS   fs.FS
}

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

func (c *Capsule) validate(request string) (string, error) {
	if len(request) > 1024 {
		return "", fmt.Errorf("too long")
	}
	u, err := url.Parse(strings.TrimSpace(request))
	if err != nil {
		return "", fmt.Errorf("invalid url")
	}
	if !u.IsAbs() {
		return "", fmt.Errorf("not an absolute url")
	}
	if u.Fragment != "" {
		return "", fmt.Errorf("fragments not allowed")
	}
	if u.User != nil {
		return "", fmt.Errorf("userinfo not allowed")
	}
	if u.Scheme != "gemini" {
		return "", fmt.Errorf("this is a gemini server")
	}
	return u.Path, nil
}

func (c *Capsule) Handle(request string, w io.Writer) error {
	if c.FS == nil {
		if c.Root == "" {
			c.Root = "."
		}
		c.FS = os.DirFS(c.Root)
	}

	path, err := c.validate(request)
	if err != nil {
		fmt.Fprintf(w, "%d %s\r\n", BadRequest, err.Error())
		return err
	}
	err = c.request(path, w)
	if err != nil {
		fmt.Fprintf(w, "%d %s\r\n", NotFound, err.Error())
		return err
	}

	return nil
}

func (c *Capsule) request(path string, w io.Writer) error {
	if path == "" {
		path = "/"
	}
	if path[len(path)-1] == '/' {
		path = path + "index.gmi"
	}
	path = strings.TrimPrefix(path, "/")

	mime := natto.Mime(path)
	switch mime {
	case "application/cgi":
		info, err := fs.Stat(c.FS, path)
		if err != nil {
			fmt.Fprintf(w, "%d %s\r\n", NotFound, err.Error())
			return fmt.Errorf("file not found")
		}
		return natto.Cgi(w, c.Root+"/"+info.Name(), "gemini")
	default:
		path = strings.TrimPrefix(path, "/")
		f, err := c.FS.Open(path)
		if err != nil {
			fmt.Fprintf(w, "%d %s\r\n", NotFound, err.Error())
			return fmt.Errorf("file not found")
		}
		defer f.Close()
		fmt.Fprintf(w, "%d %s\r\n", Success, mime)
		io.Copy(w, f)
	}
	return nil
}
