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

To run all tests run the following command:
```bash
make testacc
```

## Inspirations
* [hashicorp - terraform-provider-ad](https://github.com/hashicorp/terraform-provider-ad)

<!-- Badges -->
[godoc badge]: https://pkg.go.dev/badge/github.com/d-strobel/gowindows
[godoc page]: https://pkg.go.dev/github.com/d-strobel/gowindows

[goreport badge]: https://goreportcard.com/badge/github.com/d-strobel/gowindows
[goreport page]: https://goreportcard.com/report/github.com/d-strobel/gowindows

[build badge]: https://github.com/d-strobel/gowindows/actions/workflows/build.yml/badge.svg
[build page]: https://github.com/d-strobel/gowindows/actions/workflows/build.yml
