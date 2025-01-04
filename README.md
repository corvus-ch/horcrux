# Horcrux - a helper for preparing backups of data worth protecting

[![Build Status](https://travis-ci.org/corvus-ch/horcrux.svg?branch=master)](https://travis-ci.org/corvus-ch/horcrux)
[![Maintainability](https://api.codeclimate.com/v1/badges/58cc94f18c45c113f769/maintainability)](https://codeclimate.com/github/corvus-ch/horcrux/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/58cc94f18c45c113f769/test_coverage)](https://codeclimate.com/github/corvus-ch/horcrux/test_coverage)
[![DEB Repository](https://img.shields.io/badge/deb-packagecloud.io-844fec.svg)](https://packagecloud.io/corvus-ch/tools)
[![RPM Repository](https://img.shields.io/badge/rpm-packagecloud.io-844fec.svg)](https://packagecloud.io/corvus-ch/tools)

The use case which brought this tool into existence, is that of creating
backups of GPG private keys or any other cryptographic key or data of similar
value for that matters. Inspired by [paperkey], the idea is to transform the
data into printable form in such a way that easy recovery into its digital
original is possible. Other than [paperkey], this tool allow to split the data
into fragments, which on its own are worthless and only if brought together,
can recreate the secret they were created from.  

## Install

### deb

    curl -s https://packagecloud.io/install/repositories/corvus-ch/tools/script.deb.sh | sudo bash
    apt install bilocation

### rpm

    curl -s https://packagecloud.io/install/repositories/corvus-ch/tools/script.rpm.sh | sudo bash
    yum install bilocation

### Homebrew

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

Create PDF output using a template

Contents of `text.tmpl` located in the current working directory.

```go-tmpl
{{define "header" -}}
= My document title
{{ printf "%03d" .Output.X }}, {docdate}
:version-label: Fragment
:doctype: book

== File Information

[horizontal]
Name:: {{ .Input.Name }}
Size:: {{ .Input.Size }} bytes
MD5:: {{ .Input.Checksums.Md5 }}
SHA1:: {{ .Input.Checksums.Sha1 }}
SHA256:: {{ .Input.Checksums.Sha256 }}
SHA512:: {{ .Input.Checksums.Sha512 }}

== Text Data
....
{{end}}

{{define "footer" -}}
....

== QR Codes

image::{{ .Input.Stem }}.{{ printf "%03d" .Output.X }}.1.png[align="center"]
{{end}}
```

    horcrux create --encrypt -f text -f qr "/path/to/secret"
    for f in test.txt.*; do asciidoctor-pdf -o "${f%.txt.*}-${f##*.}.pdf" "${f}"; done

Print the PDFs. Write the corresponding password onto each of them by hand.

IMPORTANT: This produces an incomplete result, should your secret does not fit in one single QR code.
Add more image lines as needed.

## Known issues

The QR code format provide a limited feature set and can not be used to recover
the data directly. A tool like `zbarimg` from the [zbar libary][zbar] can be used to
scan the qr codes so it can be read by the zbase32 format.

## Milestones

* [x] Basic application
* [x] Plain text format for print and easy scan/ocr
* [x] QR Code format for easier scanning
* [x] Template system for custom output
* [ ] Extend possibilities with templates

## Contributing and license

This library is licenced under [MIT]. For information about how to contribute
to this project, see [CONTRIBUTING.md].

[CONTRIBUTING.md]: https://github.com/corvus-ch/horcrux/blob/master/CONTRIBUTING.md
[MIT]: https://github.com/corvus-ch/horcrux/blob/master/LICENSE
[paperkey]: http://www.jabberwocky.com/software/paperkey/
[zbar]: https://zbar.sourceforge.io
