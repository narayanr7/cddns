build:
	go build -o bin/cddns . 

build-arm:
	GOOS=linux GOARCH=arm go build -o bin/cddns . 
build-amd64:
	GOOS=linux GOARCH=amd64 go build -o bin/cddns . 
install:
	./scripts/install.sh

uninstall:
	./scripts/uninstall.sh

clean:
	rm bin/*
