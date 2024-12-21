all: natto daizu

again: clean all

daizu: natto.go cmd/daizu/main.go
	go build -o daizu cmd/daizu/main.go

natto: natto.go cmd/natto/main.go
	go build -o natto cmd/natto/main.go
	
clean:
	rm -f natto daizu

test:
	go test -cover

install:
	install natto /usr/local/bin
	install daizu /usr/local/bin

push:
	got send
	git push github

fmt:
	gofmt -s -w *.go cmd/*/main.go

README.md: README.gmi
	sisyphus -f markdown <README.gmi >README.md

doc: README.md

release: push
	git push github --tags

