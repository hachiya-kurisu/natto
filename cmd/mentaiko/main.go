package main

import (
	"blekksprut.net/natto/spartan"
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
		if !strings.HasPrefix(u, "spartan://") {
			u = "spartan://" + u
		}

		res, err := spartan.Request(ctx, u, spartan.Data{})
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		defer res.Close()

		switch res.Status {
		case spartan.ClientError, spartan.ServerError:
			fmt.Fprintln(os.Stderr, res.Status, res.Header)
		case spartan.Success:
			if *s {
				fmt.Fprintln(os.Stderr, res.Status, res.Header)
			}
			io.Copy(os.Stdout, res)
		}
	}
}
