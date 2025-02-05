all: build test

build:
	go build -o tests/find ex00/find.go
	go build -o tests/wc ex01/wc.go
	go build -o tests/xargs ex02/xargs.go
	go build -o tests/rotate ex03/rotate.go
	mkdir bin
	cp tests/find bin/find
	cp tests/wc bin/wc
	cp tests/xargs bin/xargs
	cp tests/rotate bin/rotate

test:
	go test ./tests -v
	rm tests/find tests/rotate tests/wc tests/xargs

clean:
	rm -rf bin tests/find tests/rotate tests/wc tests/xargs