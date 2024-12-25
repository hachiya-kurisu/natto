package main

import (
	"blekksprut.net/natto"
	"blekksprut.net/natto/gemini"
	"blekksprut.net/natto/spartan"
	"bufio"
	"flag"
	"log"
	"os"
)

func main() {
	r := flag.String("r", "/var/gemini", "root directory")
	s := flag.Bool("s", false, "spartan 💪")
	flag.Parse()

	err := os.Chdir(*r)
	if err != nil {
		log.Fatal("unable to chdir to root directory")
	}
	Lockdown(*r)

	var capsule natto.Capsule
	if *s {
		capsule = &gemini.Capsule{Path: *r}
	} else {
		capsule = &spartan.Capsule{Path: *r}
	}

	reader := bufio.NewReader(os.Stdin)
	request, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err.Error())
	}
	capsule.Handle(request, os.Stdout)
}
