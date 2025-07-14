BINARY=92.sh

build-armv7:
	GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=0 go build -o $(BINARY) ./cmd/crylic

build:
	go build -o $(BINARY) ./cmd/crylic

clean:
	rm -f $(BINARY)