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

package_name  = horcrux
test_files    = $(shell find . -name '*_test.go')
coverage_file = c.out

# -------
# Targets
# -------

.PHONY: install
install:
	go mod vendor

build: $(coverage_file) ${package_name}

.PHONY: test
test: $(coverage_file) test.bin ${package_name}
	./${package_name} create test.bin && ls test.txt.* > /dev/null
	./${package_name} restore -o result.bin test.txt.* && diff test.bin result.bin > /dev/null

$(coverage_file): $(test_files)
	go test -covermode=atomic -coverprofile=$@ ./...

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
	rm -f "${package_name}" test.bin test.txt.* $(coverage_file)
