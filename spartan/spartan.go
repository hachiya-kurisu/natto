package spartan

import (
	"blekksprut.net/natto"
	"bufio"
	"context"
	"fmt"
	"io"
	"io/fs"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
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

type Response struct {
	Raw    io.Reader
	Conn   net.Conn
	Status string
	Header string
}

func (r *Response) Close() {
	r.Conn.Close()
}

func (r *Response) Read(b []byte) (int, error) {
	return r.Raw.Read(b)
}

func Request(ctx context.Context, rawURL string) (*Response, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	if u.Port() == "" {
		u.Host = u.Host + ":300"
	}
	timeout, _ := time.ParseDuration("30s")
	dialer := net.Dialer{Timeout: timeout}
	conn, err := dialer.DialContext(ctx, "tcp", u.Host)
	if err != nil {
		return nil, err
	}

	if u.Path == "" {
		u.Path = "/"
	}

	fmt.Fprintf(conn, "%s %s %d\r\n", u.Hostname(), u.Path, 0)

	r := bufio.NewReader(conn)
	header, err := r.ReadString('\n')
	if err != nil {
		return nil, err
	}
	status, header, _ := strings.Cut(header, " ")
	if len(status) != 1 {
		return nil, fmt.Errorf("malformed status code")
	}
	i, err := strconv.Atoi(status)
	if err != nil || i < 2 || i > 5 {
		return nil, fmt.Errorf("invalid status code", status)
	}
	return &Response{Raw: r, Conn: conn, Status: status, Header: header}, nil
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
