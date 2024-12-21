package main

import (
	"blekksprut.net/natto"
	"flag"
	"os"
	"syscall"
)

func main() {
	r := flag.String("r", "/var/gemini", "root directory")
	c := flag.Bool("c", true, "chroot to root directory")
	s := flag.Bool("s", false, "spartan 💪")

	flag.Parse()

	capsule := natto.Capsule{Path: *r, Writer: os.Stdout}

	if *s {
		capsule.Protocol = natto.Spartan
	}

	if *c {
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

	var host, path string
	var err error

	switch capsule.Protocol {
	case natto.Spartan:
		host, path, err = capsule.SpartanRequest(os.Stdin)
		if err != nil {
			capsule.Panic(5, err.Error())
			os.Exit(1)
		}
	case natto.Gemini:
		host, path, err = capsule.GeminiRequest(os.Stdin)
		if err != nil {
			capsule.Panic(59, err.Error())
			os.Exit(1)
		}
	}

	capsule.Request(host, path)
}

