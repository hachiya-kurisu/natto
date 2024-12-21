package natto

import (
	"net/url"
	"testing"
)

var c Capsule

func TestGemtext(t *testing.T) {
	url, _ := url.Parse("gemini://test/README.gmi")
	c.Validate(url)
	if c.Request(url) != Ok {
		t.Errorf("request to this file failed")
	}
}

func TestDefaultMime(t *testing.T) {
	url, _ := url.Parse("gemini://test/natto_test.go")
	c.Validate(url)
	if c.Request(url) != Ok {
		t.Errorf("request to this file failed")
	}
}

func TestMissingFile(t *testing.T) {
	url, _ := url.Parse("gemini://test/eyyyyy")
	c.Validate(url)
	if c.Request(url) != Oops {
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
	if c.Request(url) != Oops {
		t.Errorf("shouldn't panic, at least")
	}
}
