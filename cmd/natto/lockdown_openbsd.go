package main

import (
	"golang.org/x/sys/unix"
)

func Lockdown(path string) {
	unix.Unveil(path, "r")
	unix.UnveilBlock()
	unix.Pledge("stdio exec rpath wpath", "")
}

