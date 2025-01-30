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

type Space struct {
	Root string
	FS   fs.FS
}

type Response struct {
	URL    *url.URL
	Raw    io.Reader
	Conn   net.Conn
	Status int
	Header string
}

type Data struct {
	Length int64
	Data   io.Reader
}

func (r *Response) Close() {
	r.Conn.Close()
}

func (r *Response) Read(b []byte) (int, error) {
	return r.Raw.Read(b)
}

func Request(ctx context.Context, rawURL string, d Data) (*Response, error) {
	return req(ctx, rawURL, d, 0)
}

func req(ctx context.Context, rawURL string, d Data, n int) (*Response, error) {
	if n > 5 {
		return nil, fmt.Errorf("too many redirects")
	}
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

	fmt.Fprintf(conn, "%s %s %d\r\n", u.Hostname(), u.Path, d.Length)
	if d.Length > 0 {
		written, err := io.Copy(conn, d.Data)
		if err != nil {
			return nil, err
		}
		if written != d.Length {
			return nil, fmt.Errorf("invalid data length")
		}
	}
	r := bufio.NewReader(conn)
	header, err := r.ReadString('\n')
	if err != nil {
		return nil, err
	}
	status, header, _ := strings.Cut(header, " ")
	i, _ := strconv.Atoi(status)
	switch i {
	case Redirect:
		loc, err := url.Parse(strings.TrimSpace(header))
		if err != nil {
			return nil, fmt.Errorf("invalid redirect %s", err)
		}
		if loc.Hostname() != "" {
			return nil, fmt.Errorf("no cross-site redirects")
		}
		loc.Host = u.Hostname()
		return req(ctx, loc.String(), Data{}, n+1)
	case Success, ClientError, ServerError:
		u.Host = strings.TrimSuffix(u.Host, ":300")
		return &Response{u, r, conn, i, header}, nil
	default:
		return nil, fmt.Errorf("invalid status code")
	}
}

func (c *Space) validate(request string) (string, string, error) {
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

func (c *Space) Handle(request string, rw io.ReadWriter) error {
	if c.FS == nil {
		if c.Root == "" {
			c.Root = "."
		}
		c.FS = os.DirFS(c.Root)
	}

	path, length, err := c.validate(request)
	if err != nil {
		fmt.Fprintf(rw, "%d %s\r\n", ClientError, "invalid request")
		return err
	}
	path = strings.TrimPrefix(path, "/")

	info, err := fs.Stat(c.FS, path)
	if err != nil {
		fmt.Fprintf(rw, "%d %s\r\n", ClientError, "not found")
		return err
	}

	if info.IsDir() {
		u := fmt.Sprintf("/%s/", path)
		fmt.Fprintf(rw, "%d %s\r\n", Redirect, u)
		return err
	}

	mime := natto.Mime(path)
	switch mime {
	case "application/cgi":
		os.Setenv("CONTENT_LENGTH", length)
		return natto.Cgi(rw, c.Root+"/"+info.Name(), "spartan")
	default:
		f, err := c.FS.Open(path)
		if err != nil {
			f, err = c.FS.Open(path + ".gmi")
			if err != nil {
				fmt.Fprintf(rw, "%d %s\r\n", ServerError, "unreadable")
				return fmt.Errorf("file not found")
			}
			mime = "text/gemini"
		}
		defer f.Close()
		fmt.Fprintf(rw, "%d %s\r\n", Success, mime)
		io.Copy(rw, f)
	}
	return nil
}
