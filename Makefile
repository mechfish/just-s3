.PHONY: test clean

build: dist/just-s3

dist/just-s3: *.go
	go build -o dist/just-s3 .

test:
	go test ./...

clean:
	rm dist/*
