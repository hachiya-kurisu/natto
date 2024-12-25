all: natto

again: clean all

natto: natto.go cmd/natto/main.go
	go build -C cmd/natto -o ../../natto
	
clean:
	rm -f natto

test:
	go test -cover -coverpkg ./gemini,./spartan

cover:
	go test -coverpkg ./gemini,./spartan -coverprofile=cover.out
	go tool cover -html cover.out

cert:
	openssl genrsa -out /etc/ssl/private/gemini.key 2048
	openssl req -new -key /etc/ssl/private/gemini.key \
		-out /etc/ssl/gemini.csr
	openssl x509 -req -days 2500000 \
		-in /etc/ssl/gemini.csr -signkey /etc/ssl/private/gemini.key \
		-out /etc/ssl/gemini.crt

install:
	install natto /usr/local/bin

push:
	got send
	git push github

fmt:
	gofmt -s -w *.go */*.go cmd/*/main.go

README.md: README.gmi
	sisyphus -f markdown <README.gmi >README.md

doc: README.md

release: push
	git push github --tags

