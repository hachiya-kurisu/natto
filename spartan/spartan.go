package spartan

import (
	"blekksprut.net/natto"
	"fmt"
	"io"
	"os"
	"path/filepath"
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
	Path string
}

func (c *Capsule) validate(request string) (string, int, error) {
	request = strings.TrimSpace(request)
	components := strings.SplitN(request, " ", 3)
	if len(components) != 3 {
		return "", 0, fmt.Errorf("malformed request")
	}
	path, contentLength := components[1], components[2]
	if path[0] != '/' {
		return "", 0, fmt.Errorf("missing /")
	}
	length, err := strconv.Atoi(contentLength)
	if err != nil {
		return "", 0, fmt.Errorf("invalid content length")
	}
	return path, length, nil
}

func (c *Capsule) Handle(request string, w io.Writer) error {
	path, l, err := c.validate(request)
	if err != nil {
		fmt.Fprintf(w, "%d %s\r\n", ClientError, err.Error())
		return err
	}
	if l > 0 {
		fmt.Fprintf(w, "%d %s\r\n", ServerError, err.Error())
		return err
	}

	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("file not found")
	}

	mime := natto.Types[filepath.Ext(path)]
	if mime == "" {
		mime = "application/octet-stream"
	}

	fmt.Fprintf(w, "%d %s\r\n", Success, mime)
	io.Copy(w, f)

	return nil
}
