package natto

import (
	"net/url"
	"fmt"
	"strings"
	"testing"
)

var c Capsule
var spartan Capsule

func init() {
	spartan.Protocol = Spartan
}

func TestGemtext(t *testing.T) {
	url, _ := url.Parse("gemini://test/README.gmi")
	c.Validate(url)
	if c.Request(url.Host, url.Path) != nil {
		t.Errorf("request to this file failed")
	}
}

func TestProtocol(t *testing.T) {
	var capsule Capsule
	if capsule.ProtocolName() != "gemini" {
		t.Errorf("the zero value should be the gemini protocol")
	}
}

func TestDefaultMime(t *testing.T) {
	url, _ := url.Parse("gemini://test/natto_test.go")
	c.Validate(url)
	if c.Request(url.Host, url.Path) != nil {
		t.Errorf("request to this file failed")
	}
}

func TestMissingFile(t *testing.T) {
	url, _ := url.Parse("gemini://test/eyyyyy")
	c.Validate(url)
	if c.Request(url.Host, url.Path) == nil {
		t.Errorf("this file shouldn't even be here today")
	}
}

func TestRelativeUrl(t *testing.T) {
	url, _ := url.Parse("/test/eyyyyy")
	err := c.Validate(url)
	if err == nil {
		t.Errorf("no relative urls allowed")
	}
}

func TestFragmentInUrl(t *testing.T) {
	url, _ := url.Parse("gemini://test/eyyyyy#ok")
	err := c.Validate(url)
	if err == nil {
		t.Errorf("no fragments")
	}
}

func TestScheme(t *testing.T) {
	url, _ := url.Parse("spartan://test/eyyyyy")
	err := c.Validate(url)
	if err == nil {
		t.Errorf("only gemini scheme allowed")
	}
}

func TestUserInfoInUrl(t *testing.T) {
	url, _ := url.Parse("gemini://me@test/eyyyyy")
	err := c.Validate(url)
	if err == nil {
		t.Errorf("no userinfo. please")
	}
}

// todo: actually test stuff here
func TestNoLeadingSlash(t *testing.T) {
	url, _ := url.Parse("gemini://test")
	c.Validate(url)
	if c.Request(url.Host, url.Path) == nil {
		t.Errorf("shouldn't panic, at least")
	}
}

func TestSpartan(t *testing.T) {
	url, _ := url.Parse("spartan://test/README.gmi")
	if spartan.Request(url.Host, url.Path) != nil {
		t.Errorf("couldn't serve README.gmi")
	}
}

func TestSpartanNotFound(t *testing.T) {
	url, _ := url.Parse("spartan://test/notfound.gmi")
	if spartan.Request(url.Host, url.Path) == nil {
		t.Errorf("this file shouldn't even be here today")
	}
}

func TestGeminiRequest(t *testing.T) {
	host, _, _ := c.GeminiRequest("gemini://test")
	if host != "test" {
		t.Errorf("failed to parse gemini request")
	}
}

func TestGeminiMalformedUrl(t *testing.T) {
	_, _, err := c.GeminiRequest("\b")
	if err == nil {
		t.Errorf("shouldn't have accepted the malformed request")
	}
}

func TestGeminiInvalidRequest(t *testing.T) {
	_, _, err := c.GeminiRequest("gemini://me@lol/foo")
	if err == nil {
		t.Errorf("shouldn't have accepted the malformed request")
	}
}

func TestGeminiLongRequest(t *testing.T) {
	long := strings.Repeat("a", 1016) // 1024 - length of gemini://
	_, _, err := c.GeminiRequest("gemini://" + long)
	if err == nil {
		t.Errorf("shouldn't have accepted long request")
	}
}

func TestGeminiCgi(t *testing.T) {
	spartan.CgiHandler = func(argv0 string, argv []string, envv []string) error {
		return nil
	}
	url, _ := url.Parse("spartan://test/test.sh")
	if spartan.Request(url.Host, url.Path) != nil {
		t.Errorf("cgi should succeed")
	}
}

func TestGeminiCgiFailure(t *testing.T) {
	spartan.CgiHandler = func(argv0 string, argv []string, envv []string) error {
		return fmt.Errorf("oops")
	}
	url, _ := url.Parse("spartan://test/test.sh")
	if spartan.Request(url.Host, url.Path) == nil {
		t.Errorf("cgi should fail")
	}
}

func TestSpartanRequest(t *testing.T) {
	host, _, _ := spartan.SpartanRequest("test / 0")
	if host != "test" {
		t.Errorf("failed to parse spartan request")
	}
}

func TestSpartanRequestMissingLength(t *testing.T) {
	_, _, err := spartan.SpartanRequest("test /")
	if err == nil {
		t.Errorf("requests without a length should fail")
	}
}

func TestSpartanRequestInvalidContentLength(t *testing.T) {
	_, _, err := spartan.SpartanRequest("test / .")
	if err == nil {
		t.Errorf("requests with an invalid length should fail")
	}
}

func TestSpartanRequestMalformedPath(t *testing.T) {
	_, _, err := spartan.SpartanRequest("test : 0")
	if err == nil {
		t.Errorf("requests without a leading / should fail")
	}
}

