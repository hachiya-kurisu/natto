package main

import (
	"golang.org/x/sys/unix"
)

func Lockdown(path string) {
	unix.Pledge("stdio exec rpath wpath", "")
	unix.Unveil(path, "r")
	unix.UnveilBlock()
}

