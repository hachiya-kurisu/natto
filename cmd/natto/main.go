package main

import (
	"blekksprut.net/natto"
	"bufio"
	"flag"
	"fmt"
	"net/url"
	"os"
	"syscall"
)

func main() {
	root := flag.String("r", "/var/gemini", "root directory")
	chroot := flag.Bool("c", true, "chroot to root directory")

	flag.Parse()

	capsule := natto.Capsule{Path: *root, Writer: os.Stdout}

	if *chroot {
		err := syscall.Chroot(capsule.Path)
		if err != nil {
			panic("unable to chroot to root directory")
		}
		os.Chdir("/")
	} else {
		err := os.Chdir(capsule.Path)
		if err != nil {
			panic("unable to chdir to root directory")
		}
	}

	r := bufio.NewReaderSize(os.Stdin, 1024)
	request, tooLong, err := r.ReadLine()
	if tooLong {
		capsule.Panic(59, "request is too long")
		os.Exit(1)
	}
	if err != nil {
		capsule.Panic(40, "something went wrong...")
		os.Exit(1)
	}

	u, err := url.Parse(string(request))
	if err != nil {
		capsule.Panic(40, "trouble parsing the url...")
		os.Exit(1)
	}

	capsule.Validate(u)
	capsule.Request(u)
}
