<!-- Badges -->
[![Build][build badge]][build page]
[![GoDoc][godoc badge]][godoc page]
[![GoReport][goreport badge]][goreport page]

# gowindows

Go library to configure Windows based systems.

This package mainly focuses on providing the neccessary functions for the [terraform-provider-windows](https://github.com/d-strobel/terraform-provider-windows).

## Development

### Conventional Commits

Commit messages must follow the conventional commit guidelines.<br>
For further information, see [conventionalcommits.org](https://www.conventionalcommits.org/).

### Testing

### Unit tests
The unit tests can run without external systems. All functions that relies on external systems are mocked.

Just run the following command:<br>
It will download the go dependencies and run the unit tests.
```bash
make test
```

### Acceptance test
This is not present yet!<br>
In the future I would like to spawn a system with all needed requirements to do an integration test.

To run the acceptance test run the following command:
```bash
make testacc
```

## Third-Party libraries
* [masterzen/winrm](https://github.com/masterzen/winrm)

## Inspirations
* [hashicorp - terraform-provider-ad](https://github.com/hashicorp/terraform-provider-ad)

<!-- Badges -->
[godoc badge]: https://pkg.go.dev/badge/github.com/d-strobel/gowindows
[godoc page]: https://pkg.go.dev/github.com/d-strobel/gowindows

[goreport badge]: https://goreportcard.com/badge/github.com/d-strobel/gowindows
[goreport page]: https://goreportcard.com/report/github.com/d-strobel/gowindows

[build badge]: https://github.com/d-strobel/gowindows/actions/workflows/build.yml/badge.svg
[build page]: https://github.com/d-strobel/gowindows/actions/workflows/build.yml
