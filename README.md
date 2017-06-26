# visa

Go-based Virtual Instrument Software Architecture (VISA) resource
manager.

[![GoDoc][godoc image]][godoc link]
[![License Badge][license image]][LICENSE.txt]

## Background

The primary purpose of this package is to parse a VISA resource address
string in order to create a new VISA resource, which abstracts the
interface type---USBTMC, TCPIP, ASRL. By registering a driver for the
interface type, application developers can decide which interface types
to support and which to exclude, so as to not unnecessarily bloat their
packages.

The primary source of information is the *VPP-4.3: The VISA Library*
dated June 19, 2014, which can be found at the [IVI Specifications
webpage][ivi-specs].

## Usage

In order to not bloat an end developer's application, the desired HW interface
drivers have to be registered, similar to Go's SQL package. Currently,
there are TCPIP and USBTMC drivers available.

```go
"github.com/gotmc/visa"
_ "github.com/gotmc/visa/drivers/tcpip"
_ "github.com/gotmc/visa/drivers/usbtmc"
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

[visa][] is developed using [Scott Chacon][]'s [GitHub Flow][]. To
contribute, fork [visa][], create a feature branch, and then
submit a [pull request][].  [GitHub Flow][] is summarized as:

- Anything in the `master` branch is deployable
- To work on something new, create a descriptively named branch off of
  `master` (e.g., `new-oauth2-scopes`)
- Commit to that branch locally and regularly push your work to the same
  named branch on the server
- When you need feedback or help, or you think the branch is ready for
  merging, open a [pull request][].
- After someone else has reviewed and signed off on the feature, you can
  merge it into master.
- Once it is merged and pushed to `master`, you can and *should* deploy
  immediately.

## Testing

Prior to submitting a [pull request][], please run:

```bash
$ gofmt
$ golint
$ go vet
$ go test
```

To update and view the test coverage report:

```bash
$ go test -coverprofile coverage.out
$ go tool cover -html coverage.out
```

## License

[visa][] is released under the MIT license. Please see the
[LICENSE.txt][] file for more information.

[GitHub Flow]: http://scottchacon.com/2011/08/31/github-flow.html
[godoc image]: https://godoc.org/github.com/gotmc/visa?status.svg
[godoc link]: https://godoc.org/github.com/gotmc/visa
[ivi-specs]: http://www.ivifoundation.org/specifications/
[LICENSE.txt]: https://github.com/gotmc/visa/blob/master/LICENSE.txt
[license image]: https://img.shields.io/badge/license-MIT-blue.svg
[pull request]: https://help.github.com/articles/using-pull-requests
[Scott Chacon]: http://scottchacon.com/about.html
[visa]: https://github.com/gotmc/visa
