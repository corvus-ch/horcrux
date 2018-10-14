# How to contribute

First off, thank you for considering contributing to Horcrux. It's people like
you that make Horcrux a tool usable more than just a few.

Following these guidelines helps to communicate that you respect the time of
the developers managing and developing this open source project. In return,
they should reciprocate that respect in addressing your issue, assessing
changes, and helping you finalize your pull requests.

Horcrux is an open source project maintained by a single developer. As it is a
vital tool where errors can lead to data loss or compromised security. Help
will be highly appreciated and there are many ways to contribute:

* write about this tool on your personal blog or on social media
* writing/improving tutorials/documentation on the [project wiki][wiki]
* submitting bug reports and feature requests at the [issue queue][issues]
* submitting pull requests
* checking the code base and the external dependencies it is using for
  conceptional and security flaws

Please, do not use the issue tracker for support questions. If you struggle
with using this tool, please try one of the many friendly places in the
internet first.

Please, also do not use the issue tracker to report security vulnerabilities.
Email Christian Häusler <haeusler.christian@mac.com> instead. Use the GPG key
[0xED45DBE6B88E3269][gpg] for message encryption.

Pull request — either to fix a bug or to add a new feature — will be highly
appreciated but please keep the following rules in mind:

* Your code must be test covered
* All [tests][CI] must pass
* [CodeClimate] should not find any new issues.

Notes for running the tests:

The test suite uses the binaries provided by [libgfshare]. Those tests will fail
if the those binaries are not present. In order to run those tests, please
ensure you have [libgfshare] installed and the commands `gfsplitt` and
`gfcombine` are available in your path.

[CI]: https://travis-ci.org/corvus-ch/horcrux
[CodeClimate]: https://codeclimate.com/github/corvus-ch/horcrux
[gpg]: https://pgp.mit.edu/pks/lookup?op=get&search=0xED45DBE6B88E3269
[issues]: https://github.com/corvus-ch/horcrux/issues
[libgfshare]: https://www.digital-scurf.org/software/libgfshare
[wiki]: https://github.com/corvus-ch/horcrux/wiki
