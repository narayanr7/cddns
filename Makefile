build:
	go build -o bin/cddns cmd/main.go

build-arm:
	GOOS=linux GOARCH=arm go build -o bin/cddns cmd/main.go
build-amd64:
	GOOS=linux GOARCH=amd64 go build -o bin/cddns cmd/main.go
install:
	./scripts/install.sh

uninstall:
	./scripts/uninstall.sh

clean:
	rm bin/*
