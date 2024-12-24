package main

import (
	"blekksprut.net/natto"
	"bufio"
	"flag"
	"golang.org/x/sys/unix"
	"log"
	"os"
	"syscall"
)

func main() {
	r := flag.String("r", "/var/gemini", "root directory")
	c := flag.Bool("c", true, "chroot to root directory")
	s := flag.Bool("s", false, "spartan 💪")
	e := flag.Bool("e", false, "execute cgi scripts")
	flag.Parse()
	capsule := natto.Capsule{Path: *r, Writer: os.Stdout}
	if *s {
		capsule.Protocol = natto.Spartan
	}
	if *e {
		capsule.CgiHandler = unix.Exec
	}

	if *c {
		err := syscall.Chroot(capsule.Path)
		if err != nil {
			log.Fatal("unable to chroot to root directory")
		}
		os.Chdir("/")
	} else {
		err := os.Chdir(capsule.Path)
		if err != nil {
			log.Fatal("unable to chdir to root directory")
		}
	}

	Pledge()

	reader := bufio.NewReader(os.Stdin)
	request, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(capsule.Panic(err))
	}

	var host, path string
	switch capsule.Protocol {
	case natto.Spartan:
		host, path, err = capsule.SpartanRequest(request)
		if err != nil {
			log.Fatal(capsule.Panic(err))
		}
	case natto.Gemini:
		host, path, err = capsule.GeminiRequest(request)
		if err != nil {
			log.Fatal(capsule.Panic(err))
		}
	}

	err = capsule.Request(host, path)
	if err != nil {
		log.Fatal(err.Error())
	}
}
