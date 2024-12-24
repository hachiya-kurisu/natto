package main

import (
	"golang.org/x/sys/unix"
)

func Pledge() {
	unix.Pledge("stdio exec rpath wpath", "")
}
