NAME:=gosnowflake
VERSION:=$(shell git describe --tags --abbrev=0)
REVISION:=$(shell git rev-parse --short HEAD)
COVFLAGS:=

## Run fmt, lint and test
all: fmt lint wss cov

include gosnowflake.mak

## Run tests
test_setup: test_teardown
	python3 ci/scripts/hang_webserver.py 12345 &

test_teardown:
	kill -9 $$(ps -ewf | grep hang_webserver | grep -v grep | awk '{print $$2}') || true

test: deps test_setup
	./ci/scripts/test_component.sh

## Run Coverage tests
cov:
	make test COVFLAGS="-coverprofile=coverage.txt -covermode=atomic"

wss: deps
	./ci/scripts/wss.sh

## Lint
lint: clint

## Format source codes
fmt: cfmt
	@for c in $$(ls cmd); do \
		(cd cmd/$$c;  make fmt); \
	done

## Install sample programs
install:
	for c in $$(ls cmd); do \
		(cd cmd/$$c;  GOBIN=$$GOPATH/bin go install $$c.go); \
	done

## Build fuzz tests
fuzz-build:
	for c in $$(ls | grep -E "fuzz-*"); do \
		(cd $$c; make fuzz-build); \
	done

## Run fuzz-dsn
fuzz-dsn:
	(cd fuzz-dsn; go-fuzz -bin=./dsn-fuzz.zip -workdir=.)

## Regenerate _easyjson.go files
easyjson:
	go get github.com/josharian/intern && \
	easyjson -lower_camel_case -output_filename query_easyjson.go query.go

.PHONY: setup deps update test lint help fuzz-dsn easyjson
