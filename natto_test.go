package natto

import (
	"net/url"
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

func TestGeminiLongRequest(t *testing.T) {
	long := strings.Repeat("a", 1016) // 1024 - length of gemini://
	_, _, err := c.GeminiRequest("gemini://" + long)
	if err == nil {
		t.Errorf("shouldn't have accepted long request")
	}
}
