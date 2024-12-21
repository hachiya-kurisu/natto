all: natto

again: clean all

natto: natto.go cmd/natto/main.go
	go build -o natto cmd/natto/main.go
	
clean:
	rm -f natto

test:
	go test -cover

install:
	install natto /usr/local/bin

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

