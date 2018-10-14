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
test: c.out test.bin ${package_name}
	./${package_name} create test.bin && ls test.raw.* > /dev/null
	./${package_name} restore -o result.bin test.raw.* && diff test.bin result.bin > /dev/null

c.out: main.cov $(addsuffix pkg.cov,${test_pkgs})
	find . -name '*.cov' -exec gocoverutil -coverprofile=$@ merge {} +

%.cov:
	go test -coverprofile $@ -covermode atomic ${TEST_ARGS} ./$(@D)

${package_name}: **/*.go *.go
	go build

.PHONY: release
release:
	goreleaser --rm-dist

.PHONY: snapshot
snapshot:
	goreleaser --rm-dist --snapshot

test.bin:
	dd bs=16 count=8 if=/dev/urandom of="${@}"

.PHONY: clean
clean:
	rm -f "${package_name}"
	rm -f test.bin
	find . -name '*.cov' -delete -or -name 'c.out' -delete
	find . -maxdepth 1 -name '*.bin.*' -or -name '*.raw.*' -or -name '*.txt.*' -or -name '*.zbase32.*' -or -name 'result*' | xargs rm
