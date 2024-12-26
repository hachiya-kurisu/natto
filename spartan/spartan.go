package spartan

import (
	"blekksprut.net/natto"
	"fmt"
	"io"
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
	Path string
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
	path, length, err := c.validate(request)
	if err != nil {
		fmt.Fprintf(w, "%d %s\r\n", ClientError, "invalid request")
		return err
	}
	path = "." + path

	mime := natto.Mime(path)
	switch mime {
	case "application/cgi":
		os.Setenv("CONTENT_LENGTH", length)
		return natto.Cgi(w, path, "spartan")
	default:
		f, err := os.Open(path)
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
