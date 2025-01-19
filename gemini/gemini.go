package gemini

import (
	"blekksprut.net/natto"
	"bufio"
	"context"
	"crypto/sha512"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
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

type Capsule struct {
	Root string
	FS   fs.FS
}

type Response struct {
	URL    *url.URL
	Raw    io.Reader
	Conn   *tls.Conn
	Status int
	Header string
	Cert   *x509.Certificate
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

func (r *Response) Close() {
	r.Conn.Close()
}

func (r *Response) Signature() [64]byte {
	return sha512.Sum512(r.Cert.Raw)
}

func (r *Response) SignatureBase64() string {
	signature := r.Signature()
	return base64.StdEncoding.EncodeToString(signature[:])
}

func (r *Response) Read(b []byte) (int, error) {
	return r.Raw.Read(b)
}

func Request(ctx context.Context, rawURL string) (*Response, error) {
	return doRequest(ctx, rawURL, 0)
}

func checkCertificate(cert *x509.Certificate) error {
	now := time.Now()
	if now.Before(cert.NotBefore) || now.After(cert.NotAfter) {
		return fmt.Errorf("certificate is either expired or not yet valid")
	}
	return nil
}

func doRequest(ctx context.Context, rawURL string, n int) (*Response, error) {
	if n > 5 {
		return nil, fmt.Errorf("too many redirects")
	}
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	if u.Port() == "" {
		u.Host = u.Host + ":1965"
	}
	timeout, _ := time.ParseDuration("30s")
	nd := net.Dialer{Timeout: timeout}
	config := tls.Config{InsecureSkipVerify: true}
	dialer := tls.Dialer{NetDialer: &nd, Config: &config}

	conn, err := dialer.DialContext(ctx, "tcp", u.Host)
	if err != nil {
		return nil, err
	}

	tls := conn.(*tls.Conn)
	state := tls.ConnectionState()
	cert := state.PeerCertificates[0]

	err = cert.VerifyHostname(u.Hostname())
	if err != nil {
		return nil, err
	}

	err = checkCertificate(cert)
	if err != nil {
		return nil, err
	}

	fmt.Fprintf(conn, "%s\r\n", rawURL)
	r := bufio.NewReader(conn)
	header, err := r.ReadString('\n')
	if err != nil {
		return nil, err
	}
	status, header, _ := strings.Cut(header, " ")
	if len(status) < 2 {
		return nil, fmt.Errorf("malformed header")
	}
	i, err := strconv.Atoi(status)
	if err != nil || i < 10 || i > 69 {
		return nil, fmt.Errorf("invalid status code", status)
	}

	header = strings.TrimSpace(header)

	if i >= 30 && i <= 40 {
		loc, err := u.Parse(header)
		if err != nil {
			return nil, fmt.Errorf("invalid redirect %s", err)
		}
		return doRequest(ctx, loc.String(), n+1)
	}
	u.Host = strings.TrimSuffix(u.Host, ":1965")

	return &Response{u, r, tls, i, header, cert}, nil
}
