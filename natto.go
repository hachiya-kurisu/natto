package natto

import (
	"io"
	"path/filepath"
)

const Version = "0.1.1"

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
