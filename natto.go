package natto

import (
	"io"
)

const Version = "0.0.7"

var Types = map[string]string{
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
