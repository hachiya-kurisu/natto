package main

import (
	"golang.org/x/sys/unix"
)

func Lockdown(path string) {
	unix.Unveil(path, "r w x c")
	unix.UnveilBlock()
	unix.PledgePromises("stdio exec cpath rpath wpath proc")
}
