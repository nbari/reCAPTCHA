.PHONY: all get test clean build generate

GO ?= go

all: clean build

generate:
	${GO} generate

build: generate
#	${GO} build
	env GOOS=freebsd GOARCH=amd64 go build

clean:
	@rm -rf reCAPTCHA
