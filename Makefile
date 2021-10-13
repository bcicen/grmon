NAME=grmon
VERSION=0.1

clean:
	rm -rf _build _release

build:
	mkdir _build
	go mod download
	CGO_ENABLED=0 go build -ldflags "-w" -o _build/$(NAME)
	cd _build; sha256sum * > sha256sums.txt

install:
	make build
	mv -v _build/$(NAME) $(GOPATH)/bin/$(NAME)

release:
	mkdir _release
	cp _build/* _release
	cd _release; sha256sum --quiet --check sha256sums.txt && \
	gh _release create v$(VERSION) -d -t v$(VERSION) *
