TEST = ./gemini,./spartan,.

NATTO_GEMINI_TEST_URL ?= gemini://higeki.jp
NATTO_SPARTAN_TEST_URL ?= spartan://higeki.jp

all: natto karashi negi okra mentaiko

again: clean all

natto: natto.go gemini/gemini.go spartan/spartan.go cmd/natto/main.go
	go build -C cmd/natto -o ../../natto
	
karashi: natto.go gemini/gemini.go cmd/karashi/main.go
	go build -C cmd/karashi -o ../../karashi

negi: natto.go gemini/gemini.go cmd/negi/main.go
	go build -C cmd/negi -o ../../negi

okra: natto.go gemini/gemini.go cmd/okra/main.go
	go build -C cmd/okra -o ../../okra

mentaiko: natto.go spartan/spartan.go cmd/mentaiko/main.go
	go build -C cmd/mentaiko -o ../../mentaiko

clean:
	rm -f natto karashi negi okra mentaiko

test:
	NATTO_GEMINI_TEST_URL=$(NATTO_GEMINI_TEST_URL) \
	NATTO_SPARTAN_TEST_URL=$(NATTO_SPARTAN_TEST_URL) \
	go test -cover -coverpkg $(TEST)

cover:
	NATTO_GEMINI_TEST_URL=$(NATTO_GEMINI_TEST_URL) \
	NATTO_SPARTAN_TEST_URL=$(NATTO_SPARTAN_TEST_URL) \
	go test -coverpkg $(TEST) -coverprofile=cover.out
	go tool cover -html cover.out

cert:
	openssl genrsa -out /etc/ssl/private/gemini.key 2048
	openssl req -new -key /etc/ssl/private/gemini.key \
		-out /etc/ssl/gemini.csr
	openssl x509 -req -days 2500000 \
		-in /etc/ssl/gemini.csr -signkey /etc/ssl/private/gemini.key \
		-out /etc/ssl/gemini.crt

install: all
	install natto /usr/local/bin/
	install karashi /usr/local/bin/
	install negi /usr/local/bin/
	install okra /usr/local/bin/
	install mentaiko /usr/local/bin/

push:
	got send
	git push github

fmt:
	gofmt -s -w *.go */*.go cmd/*/*.go

README.md: README.gmi
	sisyphus -f markdown <README.gmi >README.md

doc: README.md

dead:
	deadcode ./...

release: push
	git push github --tags

