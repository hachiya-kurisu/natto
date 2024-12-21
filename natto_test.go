package natto

import (
	"io/ioutil"
	"net/url"
	"testing"
)

var c Capsule

func init() {
	c.Writer = ioutil.Discard
}

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
