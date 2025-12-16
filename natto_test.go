package natto_test

import (
	"bytes"
	"context"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"blekksprut.net/natto/gemini"
	"blekksprut.net/natto/spartan"
)

var g gemini.Capsule
var s spartan.Space

func TestGemini(t *testing.T) {
	err := g.Handle("gemini://localhost/README.gmi", &bytes.Buffer{})
	if err != nil {
		t.Errorf("request shouldn't have failed")
	}
}

func TestMissingSlash(t *testing.T) {
	err := g.Handle("gemini://localhost", &bytes.Buffer{})
	if err == nil {
		t.Errorf("request shouldn't have failed")
	}
}

func TestCgi(t *testing.T) {
	err := g.Handle("gemini://localhost/hello.cgi", &bytes.Buffer{})
	if err != nil {
		t.Errorf("request shouldn't have failed")
	}
}

func TestBrokenCgi(t *testing.T) {
	err := g.Handle("gemini://localhost/failure.cgi", &bytes.Buffer{})
	if err == nil {
		t.Errorf("request should have failed")
	}
}

func TestMissingCgi(t *testing.T) {
	err := g.Handle("gemini://localhost/notfound.cgi", &bytes.Buffer{})
	if err == nil {
		t.Errorf("request should have failed")
	}
}

func TestDefaultMime(t *testing.T) {
	err := g.Handle("gemini://localhost/natto.go", &bytes.Buffer{})
	if err != nil {
		t.Errorf("request shouldn't have failed")
	}
}

func TestMissingFile(t *testing.T) {
	err := g.Handle("gemini://localhost/notFound", &bytes.Buffer{})
	if err == nil {
		t.Errorf("request should have failed")
	}
}

func TestRequestLength(t *testing.T) {
	err := g.Handle(strings.Repeat("_", 1025), &bytes.Buffer{})
	if err == nil {
		t.Errorf("request should have failed")
	}
}

func TestInvalidUrl(t *testing.T) {
	err := g.Handle("\b", &bytes.Buffer{})
	if err == nil {
		t.Errorf("request should have failed")
	}
}

func TestRelativeUrl(t *testing.T) {
	err := g.Handle("/relative", &bytes.Buffer{})
	if err == nil {
		t.Errorf("shouldn't allow relative urls")
	}
}

func TestFragment(t *testing.T) {
	err := g.Handle("gemini://localhost/README.gmi#fragment", &bytes.Buffer{})
	if err == nil {
		t.Errorf("shouldn't allow fragments")
	}
}

func TestScheme(t *testing.T) {
	err := g.Handle("spartan://localhost/README.gmi", &bytes.Buffer{})
	if err == nil {
		t.Errorf("should only allow gemini:// scheme")
	}
}

func TestUserInfo(t *testing.T) {
	err := g.Handle("gemini://user@localhost/README.gmi", &bytes.Buffer{})
	if err == nil {
		t.Errorf("shouldn't allow userinfo in request")
	}
}

func TestSpartan(t *testing.T) {
	err := s.Handle("localhost /README.gmi 0", &bytes.Buffer{})
	if err != nil {
		t.Errorf("request shouldn't have failed")
	}
}

func TestSpartanMissingFile(t *testing.T) {
	err := s.Handle("localhost /notfound 0", &bytes.Buffer{})
	if err == nil {
		t.Errorf("request should have failed")
	}
}

func TestSpartanMalformedRequest(t *testing.T) {
	err := s.Handle("localhost 0", &bytes.Buffer{})
	if err == nil {
		t.Errorf("request should have failed")
	}
}

func TestSpartanMissingSlash(t *testing.T) {
	err := s.Handle("localhost oops 0", &bytes.Buffer{})
	if err == nil {
		t.Errorf("request should have failed")
	}
}

func TestSpartanInvalidContentLength(t *testing.T) {
	err := s.Handle("localhost /README.gmi zero", &bytes.Buffer{})
	if err == nil {
		t.Errorf("request should have failed")
	}
}

func TestSpartanDefaultMime(t *testing.T) {
	err := s.Handle("localhost /natto.go 0", &bytes.Buffer{})
	if err != nil {
		t.Errorf("request shouldn't have failed")
	}
}

func TestSpartanIndex(t *testing.T) {
	err := s.Handle("localhost / 0", &bytes.Buffer{})
	if err == nil {
		t.Errorf("request should have failed")
	}
}

func TestSpartanCgi(t *testing.T) {
	err := s.Handle("localhost /hello.cgi 0", &bytes.Buffer{})
	if err != nil {
		t.Errorf("request shouldn't have failed")
	}
}

func TestSpartanMissingCgi(t *testing.T) {
	err := s.Handle("localhost /notfound.cgi 0", &bytes.Buffer{})
	if err == nil {
		t.Errorf("request should have failed")
	}
}

func TestSpartanRequest(t *testing.T) {
	u := os.Getenv("NATTO_SPARTAN_TEST_URL")
	if u != "" {
		r, err := spartan.Request(context.Background(), u, spartan.Data{})
		if err != nil {
			t.Errorf("failed to get spartan test url")
		} else {
			defer r.Close()
			ioutil.ReadAll(r)
		}
	}
}

func TestGeminiRequest(t *testing.T) {
	u := os.Getenv("NATTO_GEMINI_TEST_URL")
	if u != "" {
		r, err := gemini.Request(context.Background(), u)
		if err != nil {
			t.Errorf("failed to get spartan test url")
		} else {
			defer r.Close()
			ioutil.ReadAll(r)
		}
	}
}
