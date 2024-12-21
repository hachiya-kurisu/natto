package main

import (
	"blekksprut.net/natto"
	"bufio"
	"flag"
	"os"
	"strconv"
	"strings"
	"syscall"
)

func main() {
	root := flag.String("r", "/var/gemini", "root directory")
	chroot := flag.Bool("c", true, "chroot to root directory")

	flag.Parse()

	capsule := natto.Capsule{
		Path:     *root,
		Writer:   os.Stdout,
		Protocol: natto.Spartan,
	}

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

	r := bufio.NewReader(os.Stdin)
	request, err := r.ReadString('\n')
	if err != nil {
		capsule.Panic(4, "something went wrong...")
		os.Exit(1)
	}

	request = strings.TrimSpace(request)
	components := strings.SplitN(request, " ", 3)
	if len(components) != 3 {
		capsule.Panic(5, "invalid request: need host, path and length")
		return
	}
	host, path, contentLength := components[0], components[1], components[2]

	if path[0] != '/' {
		capsule.Panic(5, "invalid request: path has to begin with a /")
		return
	}

	length, err := strconv.Atoi(contentLength)
	if err != nil {
		capsule.Panic(5, "invalid content length")
		return
	}
	if length > 0 {
		capsule.Panic(4, "no data block support (yet)")
		return
	}

	// err = capsule.Validate(u)
	// if err != nil {
	// 	capsule.Panic(59, err.Error())
	// }

	capsule.Request(host, path)
}
