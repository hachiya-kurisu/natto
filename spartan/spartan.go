package spartan

import (
	"blekksprut.net/natto"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strconv"
	"strings"
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
	path, length := components[1], components[2]
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

	mime := natto.Mime(path)
	switch mime {
	case "application/cgi":
		os.Setenv("CONTENT_LENGTH", length)
		info, err := fs.Stat(c.FS, path)
		if err != nil {
			fmt.Fprintf(w, "%d %s\r\n", ClientError, err.Error())
			return fmt.Errorf("file not found")
		}
		return natto.Cgi(w, c.Root+"/"+info.Name(), "spartan")
	default:
		f, err := c.FS.Open(path)
		if err != nil {
			fmt.Fprintf(w, "%d %s\r\n", ClientError, "not found")
			return fmt.Errorf("file not found")
		}
		defer f.Close()
		fmt.Fprintf(w, "%d %s\r\n", Success, mime)
		io.Copy(w, f)
	}
	return nil
}
