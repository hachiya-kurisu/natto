package natto_test

import (
	"blekksprut.net/natto/gemini"
	"blekksprut.net/natto/spartan"
	"io"
	"strings"
	"testing"
)

var g gemini.Capsule
var s spartan.Capsule

func TestGemini(t *testing.T) {
	err := g.Handle("gemini://localhost/README.gmi", io.Discard)
	if err != nil {
		t.Errorf("request shouldn't have failed")
	}
}

func TestMissingSlash(t *testing.T) {
	err := g.Handle("gemini://localhost", io.Discard)
	if err == nil {
		t.Errorf("request shouldn't have failed")
	}
}

func TestCgi(t *testing.T) {
	err := g.Handle("gemini://localhost/hello.cgi", io.Discard)
	if err != nil {
		t.Errorf("request shouldn't have failed")
	}
}

func TestBrokenCgi(t *testing.T) {
	err := g.Handle("gemini://localhost/failure.cgi", io.Discard)
	if err == nil {
		t.Errorf("request should have failed")
	}
}

func TestMissingCgi(t *testing.T) {
	err := g.Handle("gemini://localhost/notfound.cgi", io.Discard)
	if err == nil {
		t.Errorf("request should have failed")
	}
}

func TestDefaultMime(t *testing.T) {
	err := g.Handle("gemini://localhost/natto.go", io.Discard)
	if err != nil {
		t.Errorf("request shouldn't have failed")
	}
}

func TestMissingFile(t *testing.T) {
	err := g.Handle("gemini://localhost/notFound", io.Discard)
	if err == nil {
		t.Errorf("request should have failed")
	}
}

func TestRequestLength(t *testing.T) {
	err := g.Handle(strings.Repeat("_", 1025), io.Discard)
	if err == nil {
		t.Errorf("request should have failed")
	}
}

func TestInvalidUrl(t *testing.T) {
	err := g.Handle("\b", io.Discard)
	if err == nil {
		t.Errorf("request should have failed")
	}
}

func TestRelativeUrl(t *testing.T) {
	err := g.Handle("/relative", io.Discard)
	if err == nil {
		t.Errorf("shouldn't allow relative urls")
	}
}

func TestFragment(t *testing.T) {
	err := g.Handle("gemini://localhost/README.gmi#fragment", io.Discard)
	if err == nil {
		t.Errorf("shouldn't allow fragments")
	}
}

func TestScheme(t *testing.T) {
	err := g.Handle("spartan://localhost/README.gmi", io.Discard)
	if err == nil {
		t.Errorf("should only allow gemini:// scheme")
	}
}

func TestUserInfo(t *testing.T) {
	err := g.Handle("gemini://user@localhost/README.gmi", io.Discard)
	if err == nil {
		t.Errorf("shouldn't allow userinfo in request")
	}
}

func TestSpartan(t *testing.T) {
	err := s.Handle("localhost /README.gmi 0", io.Discard)
	if err != nil {
		t.Errorf("request shouldn't have failed")
	}
}

func TestSpartanMissingFile(t *testing.T) {
	err := s.Handle("localhost /notfound 0", io.Discard)
	if err == nil {
		t.Errorf("request should have failed")
	}
}

func TestSpartanMalformedRequest(t *testing.T) {
	err := s.Handle("localhost 0", io.Discard)
	if err == nil {
		t.Errorf("request should have failed")
	}
}

func TestSpartanMissingSlash(t *testing.T) {
	err := s.Handle("localhost oops 0", io.Discard)
	if err == nil {
		t.Errorf("request should have failed")
	}
}

func TestSpartanInvalidContentLength(t *testing.T) {
	err := s.Handle("localhost /README.gmi zero", io.Discard)
	if err == nil {
		t.Errorf("request should have failed")
	}
}

func TestSpartanDefaultMime(t *testing.T) {
	err := s.Handle("localhost /natto.go 0", io.Discard)
	if err != nil {
		t.Errorf("request shouldn't have failed")
	}
}

func TestSpartanIndex(t *testing.T) {
	err := s.Handle("localhost / 0", io.Discard)
	if err == nil {
		t.Errorf("request should have failed")
	}
}

func TestSpartanCgi(t *testing.T) {
	err := s.Handle("localhost /hello.cgi 0", io.Discard)
	if err != nil {
		t.Errorf("request shouldn't have failed")
	}
}

func TestSpartanMissingCgi(t *testing.T) {
	err := s.Handle("localhost /notfound.cgi 0", io.Discard)
	if err == nil {
		t.Errorf("request should have failed")
	}
}
