package spartan

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"strconv"
	"strings"

	"blekksprut.net/natto"
)

const (
	Success     = 2
	Redirect    = 3
	ClientError = 4
	ServerError = 5
)

type Capsule struct {
	Root string
	FS   fs.FS
}

func (c *Capsule) validate(request string) (string, string, error) {
	request = strings.TrimSpace(request)
	components := strings.SplitN(request, " ", 3)
	if len(components) != 3 {
		return "", "", fmt.Errorf("malformed request")
	}
	_, path, length := components[0], components[1], components[2]
	if path[0] != '/' {
		return "", "", fmt.Errorf("missing /")
	}
	if path[len(path)-1] == '/' {
		path = path + "index.gmi"
	}
	_, err := strconv.Atoi(length)
	if err != nil {
		return "", "", fmt.Errorf("invalid content length")
	}
	return path, length, nil
}

func (c *Capsule) Handle(request string, w io.Writer) error {
	if c.FS == nil {
		if c.Root == "" {
			c.Root = "."
		}
		c.FS = os.DirFS(c.Root)
	}

	path, length, err := c.validate(request)
	if err != nil {
		fmt.Fprintf(w, "%d %s\r\n", ClientError, "invalid request")
		return err
	}
	path = strings.TrimPrefix(path, "/")

	info, err := fs.Stat(c.FS, path)
	if err != nil {
		fmt.Fprintf(w, "%d %s\r\n", ClientError, "not found")
		return err
	}

	if info.IsDir() {
		u := fmt.Sprintf("/%s/", path)
		fmt.Fprintf(w, "%d %s\r\n", Redirect, u)
		return err
	}

	mime := natto.Mime(path)
	switch mime {
	case "application/cgi":
		os.Setenv("CONTENT_LENGTH", length)
		return natto.Cgi(w, c.Root+"/"+info.Name(), "spartan")
	default:
		f, err := c.FS.Open(path)
		if err != nil {
			f, err = c.FS.Open(path + ".gmi")
			if err != nil {
				fmt.Fprintf(w, "%d %s\r\n", ServerError, "unreadable")
				return fmt.Errorf("file not found")
			}
			mime = "text/gemini"
		}
		defer f.Close()
		fmt.Fprintf(w, "%d %s\r\n", Success, mime)
		io.Copy(w, f)
	}
	return nil
}
