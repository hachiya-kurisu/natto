package natto

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

const Version = "0.1.7"

var Types = map[string]string{
	".cgi":  "application/cgi",
	".gmi":  "text/gemini",
	".txt":  "text/plain",
	".jpg":  "image/jpeg",
	".png":  "image/png",
	".jxl":  "image/jxl",
	".webp": "image/webp",
	".mp3":  "audio/mpeg",
	".m4a":  "audio/mp4",
	".mp4":  "video/mp4",
}

type Capsule interface {
	Handle(string, io.Writer) error
}

func Mime(path string) string {
	mime := Types[filepath.Ext(path)]
	if mime == "" {
		mime = "application/octet-stream"
	}
	return mime
}

func Cgi(w io.Writer, path string, protocol string) error {
	os.Setenv("GATEWAY_INTERFACE", "CGI/1.1")
	os.Setenv("SERVER_PROTOCOL", protocol)
	cmd := exec.Command(path)
	cmd.Env = os.Environ()
	cmd.Stdout = w
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("cgi trouble: %s", err.Error())
	}
	return nil
}
