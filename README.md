# visa

Go-based Virtual Instrument Software Architecture (VISA) resource manager.

[![GoDoc][godoc badge]][godoc link]
[![Go Report Card][report badge]][report card]
[![License Badge][license badge]][LICENSE.txt]

## Background

The primary purpose of this package is to parse a VISA resource address
string in order to create a new VISA resource, which abstracts the
interface type—USBTMC, TCPIP, ASRL. By registering a driver for the
interface type, application developers can decide which interface types
to support and which to exclude, so as to not unnecessarily bloat their
packages.

The primary source of information is the *VPP-4.3: The VISA Library*
dated June 19, 2014, which can be found at the [IVI Specifications
webpage][ivi-specs].

## Usage

In order to not bloat an end developer's application, the desired HW interface
driver(s) have to be registered, similar to Go's SQL package. Currently, there
are TCPIP, USBTMC, and Serial (ASRL) drivers available.

```go
"github.com/gotmc/visa"
_ "github.com/gotmc/visa/drivers/tcpip"
_ "github.com/gotmc/visa/drivers/usbtmc"
_ "github.com/gotmc/visa/drivers/asrl"
```

## Installation

```bash
$ go get github.com/gotmc/visa
```

## Documentation

Documentation can be found at either:

- <https://godoc.org/github.com/gotmc/visa>
- <http://localhost:6060/pkg/github.com/gotmc/visa/> after running `$
  godoc -http=:6060`

## Contributing

Contributions are welcome! To contribute please:

1. Fork the repository
2. Create a feature branch
3. Code
4. Submit a [pull request][]

### Testing

Prior to submitting a [pull request][], please run:

```bash
$ make check
$ make lint
```

To update and view the test coverage report:

```bash
$ make cover
```

## License

[visa][] is released under the MIT license. Please see the
[LICENSE.txt][] file for more information.

[GitHub Flow]: http://scottchacon.com/2011/08/31/github-flow.html
[godoc badge]: https://godoc.org/github.com/gotmc/visa?status.svg
[godoc link]: https://godoc.org/github.com/gotmc/visa
[ivi-specs]: http://www.ivifoundation.org/specifications/
[LICENSE.txt]: https://github.com/gotmc/visa/blob/master/LICENSE.txt
[license badge]: https://img.shields.io/badge/license-MIT-blue.svg
[pull request]: https://help.github.com/articles/using-pull-requests
[report badge]: https://goreportcard.com/badge/github.com/gotmc/visa
[report card]: https://goreportcard.com/report/github.com/gotmc/visa
[visa]: https://github.com/gotmc/visa
