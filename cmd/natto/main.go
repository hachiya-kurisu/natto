package main

import (
	"blekksprut.net/natto"
	"blekksprut.net/natto/gemini"
	"blekksprut.net/natto/spartan"
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	r := flag.String("r", "/var/gemini", "root directory")
	s := flag.Bool("s", false, "spartan 💪")
	v := flag.Bool("v", false, "version")
	flag.Parse()

	if *v {
		fmt.Println(os.Args[0], natto.Version)
		os.Exit(0)
	}

	path, err := filepath.Abs(*r)
	if err != nil {
		log.Fatal("invalid root path")
	}

	err = os.Chdir(path)
	if err != nil {
		log.Fatal("unable to chdir to root directory")
	}
	Lockdown(path)

	var capsule natto.Capsule
	if *s {
		capsule = &spartan.Space{Root: path}
	} else {
		capsule = &gemini.Capsule{Root: path}
	}

	reader := bufio.NewReader(os.Stdin)
	request, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err.Error())
	}
	capsule.Handle(request, &natto.Stdio{})
}
