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

## Contributing and license

This library is licenced under [MIT]. For information about how to contribute
to this project, see [CONTRIBUTING.md].

[CONTRIBUTING.md]: https://github.com/corvus-ch/horcrux/blob/master/CONTRIBUTING.md
[MIT]: https://github.com/corvus-ch/horcrux/blob/master/LICENSE
[paperkey]: http://www.jabberwocky.com/software/paperkey/
