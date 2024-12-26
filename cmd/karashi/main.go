package main

import (
	"blekksprut.net/natto"
	"blekksprut.net/natto/gemini"
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
)

func serve(socket net.Conn, capsule *gemini.Capsule) {
	defer socket.Close()
	reader := bufio.NewReader(socket)
	request, err := reader.ReadString('\n')
	if err != nil {
		log.Println(err)
		return
	}
	capsule.Handle(request, socket)
}

func main() {
	a := flag.String("a", ":1965", "address")
	c := flag.String("c", "/etc/ssl/gemini.crt", "certificate")
	k := flag.String("k", "/etc/ssl/private/gemini.key", "private key")
	r := flag.String("r", "/var/gemini", "root directory")
	v := flag.Bool("v", false, "version")

	flag.Parse()

	if *v {
		fmt.Println(os.Args[0], natto.Version)
		os.Exit(0)
	}

	cert, err := tls.LoadX509KeyPair(*c, *k)
	if err != nil {
		log.Fatal("keypair trouble")
	}
	config := tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	err = os.Chdir(*r)
	if err != nil {
		log.Fatal("unable to chdir to root directory")
	}
	Lockdown(*r)

	capsule := &gemini.Capsule{Root: *r}
	server, err := tls.Listen("tcp", *a, &config)
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
