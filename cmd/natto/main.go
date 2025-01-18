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

	err := os.Chdir(*r)
	if err != nil {
		log.Fatal("unable to chdir to root directory")
	}
	Lockdown(*r)

	var capsule natto.Capsule
	if *s {
		capsule = &spartan.Space{Root: *r}
	} else {
		capsule = &gemini.Capsule{Root: *r}
	}

	reader := bufio.NewReader(os.Stdin)
	request, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err.Error())
	}
	capsule.Handle(request, os.Stdout)
}
