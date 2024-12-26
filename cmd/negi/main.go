package main

import (
	"blekksprut.net/natto"
	"blekksprut.net/natto/gemini"
	"blekksprut.net/natto/spartan"
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
)

func serve(socket net.Conn, capsule natto.Capsule) {
	defer socket.Close()
	reader := bufio.NewReader(socket)
	request, err := reader.ReadString('\n')
	if err != nil {
		log.Println(err)
	}
	capsule.Handle(request, socket)
}

func main() {
	r := flag.String("r", "/var/gemini", "root directory")
	s := flag.Bool("s", false, "spartan 💪")
	v := flag.Bool("v", false, "version")

	flag.Parse()

	var a *string
	if *s {
		a = flag.String("a", ":300", "address")
	} else {
		a = flag.String("a", ":1965", "address")
	}

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
		capsule = &spartan.Capsule{Path: *r}
	} else {
		capsule = &gemini.Capsule{Path: *r}
	}

	server, err := net.Listen("tcp", *a)
	if err != nil {
		log.Fatal(err)
	}
	defer server.Close()

	log.Printf("listening on %s\n", *a)
	for {
		socket, err := server.Accept()
		if err != nil {
			log.Printf("unacceptable: %v", err)
		}
		go serve(socket, capsule)
	}
}
