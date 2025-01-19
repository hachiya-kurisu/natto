package main

import (
	"blekksprut.net/natto/gemini"
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	ctx := context.Background()

	s := flag.Bool("s", false, "print status line")

	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(0)
	}

	for _, u := range flag.Args() {
		if !strings.HasPrefix(u, "gemini://") {
			u = "gemini://" + u
		}

		res, err := gemini.Request(ctx, u)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		defer res.Close()

		switch res.Status[0] {
		case '1':
			fmt.Fprintf(os.Stderr, "Input requests not supported\n")
		case '6':
			fmt.Fprintf(os.Stderr, "Client certificate support not yet implemented\n")
		case '4', '5':
			fmt.Fprintf(os.Stderr, "%s %s\n", res.Status, res.Header)
		case '2':
			if *s {
				fmt.Fprintf(os.Stderr, "%s %s\n", res.Status, res.Header)
			}
			switch {
			case strings.HasPrefix(res.Header, "text/"):
				io.Copy(os.Stdout, res)
			default:
				fmt.Fprintln(os.Stderr, "only text responses supported")
			}
		default:
			fmt.Fprintf(os.Stderr, "Unknown status code %s\n", res.Status)
		}
	}
}
