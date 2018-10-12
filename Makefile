MAKEFLAGS += --warn-undefined-variables
SHELL := bash
.SHELLFLAGS := -eu -o pipefail -c
.DEFAULT_GOAL := test
.DELETE_ON_ERROR:
.SUFFIXES:

# ---------------------
# Environment variables
# ---------------------

GOPATH    ?= $(shell go env GOPATH)
TEST_ARGS ?=

# ------------------
# Internal variables
# ------------------

package_name = horcrux
test_pkgs    = $(dir $(shell find . -name '*_test.go'))

# -------
# Targets
# -------

.PHONY: install
install:
	go get -t -d -tags=integration ./...
	go get -u github.com/AlekSi/gocoverutil

build: c.out ${package_name}

.PHONY: test
test: c.out

c.out: main.cov $(addsuffix pkg.cov,${test_pkgs})
	find . -name '*.cov' -exec gocoverutil -coverprofile=$@ merge {} +

%.cov:
	go test -coverprofile $@ -covermode atomic ${TEST_ARGS} ./$(@D)

${package_name}: **/*.go *.go
	go build

.PHONY: clean
clean:
	rm -f "${package_name}"
	find . -name '*.cov' -delete -or -name 'c.out' -delete
