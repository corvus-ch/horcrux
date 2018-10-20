# Horcrux - a helper for preparing backups of data worth protecting

[![Build Status](https://travis-ci.org/corvus-ch/horcrux.svg?branch=master)](https://travis-ci.org/corvus-ch/horcrux)
[![Maintainability](https://api.codeclimate.com/v1/badges/58cc94f18c45c113f769/maintainability)](https://codeclimate.com/github/corvus-ch/horcrux/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/58cc94f18c45c113f769/test_coverage)](https://codeclimate.com/github/corvus-ch/horcrux/test_coverage)

The use case which brought this tool into existence, is that of creating
backups of GPG private keys or any other cryptographic key or data of similar
value for that matters. Inspired by [paperkey], the idea is to transform the
data into printable form in such a way that easy recovery into its digital
original is possible. Other than [paperkey], this tool allow to split the data
into fragments, which on its own are worthless and only if brought together,
can recreate the secret they were created from.  

## Install

### Linux (Debian/Ubuntu)

1. Visit https://github.com/corvus-ch/horcrux/releases and download the latest `.deb` package.
2. Use `dpkg -i` to install that package.

### OS X

    brew install corvus-ch/tools/horcrux

If you are not familiar with Homebrew visit https://brew.sh.

## Usage

Backup a GPG key:

    KEY_ID=… # Declare variable holding the ID of the GPG key you want to backup.
    gpg --export "${KEY_ID}" > public.gpg
    gpg --export-secret-key "${KEY_ID}" > "${KEY_ID}.gpg"
    paperkey --secret-key="${KEY_ID}.gpg" --output-type raw --output="${KEY_ID}.bin"
    horcrux create "${KEY_ID}.bin"
    ls *.txt.* # Those are the files you can now place at your backup locations.

Restore a GPG key (builds on top of the above example):

    horcrux restore -o paperkey.bin *.txt.* # For this example only two of the three files are required.
    paperkey --pubring=public.gpg --secrets=paperkey.bin --input-type=raw --output=secret.gpg
    diff "${KEY_ID}.gpg" secret.gpg

## Milestones

* [x] Basic application
* [x] Plain text format for print and easy scan/ocr
* [ ] QR Code format for easier scanning
* [ ] Template system for custom output

## Contributing and license

This library is licenced under [MIT]. For information about how to contribute
to this project, see [CONTRIBUTING.md].

[CONTRIBUTING.md]: https://github.com/corvus-ch/horcrux/blob/master/CONTRIBUTING.md
[MIT]: https://github.com/corvus-ch/horcrux/blob/master/LICENSE
[paperkey]: http://www.jabberwocky.com/software/paperkey/
