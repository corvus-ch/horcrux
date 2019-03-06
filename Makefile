MAKEFLAGS += --warn-undefined-variables
SHELL := bash
.SHELLFLAGS := -eu -o pipefail -c
.DEFAULT_GOAL := test
.DELETE_ON_ERROR:
.SUFFIXES:

# ---------------------
# Environment variables
# ---------------------

GOPATH     ?= $(shell go env GOPATH)
GORELEASER ?= goreleaser

# ------------------
# Internal variables
# ------------------

binary_name   = horcrux
cover_file    = c.out

# -------
# Targets
# -------

.PHONY: deps
deps:
	go mod download

.PHONY: test
test: $(cover_file) test.bin $(package_name)
test: $(coverage_file) test.bin $(package_name)
	./${package_name} create test.bin && ls test.txt.* > /dev/null
	./${package_name} restore -o result.bin test.txt.* && diff test.bin result.bin > /dev/null


$(cover_file): $(wildcard **/*_test.go)
	go test -covermode=atomic -coverprofile=$@ ./...

.PHONY: build
build: $(binary_name)

$(binary_name): $(wildcard **/*.go)
	go build

.PHONY: release
release:
	$(GORELEASER) --rm-dist
	package_cloud push corvus-ch/tools/ubuntu/xenial dist/horcrux_*_linux_*.deb
	package_cloud push corvus-ch/tools/ubuntu/bionic dist/horcrux_*_linux_*.deb
	package_cloud push corvus-ch/tools/debian/stretch dist/horcrux_*_linux_armv6.deb
	package_cloud push corvus-ch/tools/debian/buster dist/horcrux_*_linux_*.deb
	package_cloud push corvus-ch/tools/raspbian/stretch dist/horcrux_*_linux_armv6.deb
	package_cloud push corvus-ch/tools/raspbian/buster dist/horcrux_*_linux_armv6.deb
	package_cloud push corvus-ch/tools/el/6 dist/horcrux_*_linux_*.rpm
	package_cloud push corvus-ch/tools/el/7 dist/horcrux_*_linux_*.rpm
	package_cloud push corvus-ch/tools/fedora/28 dist/horcrux_*_linux_*.rpm
	package_cloud push corvus-ch/tools/fedora/29 dist/horcrux_*_linux_*.rpm

.PHONY: snapshot
snapshot:
	$(GORELEASER) --rm-dist --snapshot

test.bin:
	dd bs=16 count=8 if=/dev/urandom of="${@}"

.PHONY: clean
clean:
	rm -rf $(binary_name) $(cover_file) test.txt.* result.bin
