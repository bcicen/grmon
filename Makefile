NAME=grmon

build:
	dep ensure
	cd cmd/grmon; CGO_ENABLED=0 go build -ldflags "-w" -o ../../$(NAME)

install:
	make build
	mv -v $(NAME) $(GOPATH)/bin/$(NAME)
