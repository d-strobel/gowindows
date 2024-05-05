# gowindows
<!-- Badges -->
[![Build][build badge]][build page]
[![GoDoc][godoc badge]][godoc page]
[![GoReport][goreport badge]][goreport page]
[![Conventional Commits][convention badge]][convention page]

![logo](images/logo/gowindows-icon_250.png)

## Overview
**gowindows** is a Go library designed for remotely configuring and managing Windows-based systems.

Leveraging WinRM and SSH connections, gowindows provides a comprehensive set of functions to execute PowerShell commands, making it easy to automate tasks, manage users, groups, and more on remote Windows servers.

This library is especially useful when combined with tools like Terraform, enabling seamless integration into infrastructure as code workflows for Windows environments.

## Usage

### Single Client with an SSH Connection
```go
package main

import (
	"context"
	"fmt"

	"github.com/d-strobel/gowindows/connection/ssh"
	"github.com/d-strobel/gowindows/windows/localaccounts"
)

func main() {
	sshConfig := &ssh.Config{
		Host:     "winsrv",
		Username: "vagrant",
		Password: "vagrant",
	}

	// Create a new connection.
	conn, err := ssh.NewConnection(sshConfig)
	if err != nil {
		panic(err)
	}

	// Create a client for the localaccounts package.
	c := localaccounts.NewClient(conn)
	defer c.Connection.Close()

	// Run the GroupRead function to retrieve a local Windows group.
	group, err := c.GroupRead(context.Background(), localaccounts.GroupReadParams{Name: "Users"})
	if err != nil {
		panic(err)
	}

	// Print the user group.
	fmt.Printf("User group: %+v", group)
}
```

### Multi Client with a WinRM Connection
```go
package main

import (
	"context"
	"fmt"

	"github.com/d-strobel/gowindows"
	"github.com/d-strobel/gowindows/connection/winrm"
	"github.com/d-strobel/gowindows/windows/localaccounts"
)

func main() {
	winrmConfig := &winrm.Config{
		Host:     "winsrv",
		Username: "vagrant",
		Password: "vagrant",
	}

	// Create a new connection.
	conn, err := winrm.NewConnection(winrmConfig)
	if err != nil {
		panic(err)
	}

	// Create client for all subpackages.
	c := gowindows.NewClient(conn)
	defer c.Close()

	// Run the GroupRead function to retrieve a local Windows group.
	group, err := c.LocalAccounts.GroupRead(context.Background(), localaccounts.GroupReadParams{Name: "Users"})
	if err != nil {
		panic(err)
	}

	// Print the user group.
	fmt.Printf("User group: %+v", group)
}
```

## Development
### Conventional Commits
**gowindows** follows the conventional commit guidelines. For more information, see [conventionalcommits.org](https://www.conventionalcommits.org/).

### Testing
### Unit tests
Run unit tests:
```bash
make test
```

### Acceptance test
Prerequisites:
* [Hashicorp Vagrant](https://www.vagrantup.com/)
* [Oracle VirtualBox](https://www.virtualbox.org/)

Boot the Vagrant machines:
```bash
make vagrant-up
```

Run acceptance tests:
```bash
make testacc
```

Destroy the Vagrant machines:
```bash
make vagrant-down
```

## Third-Party libraries
* For this project, I made a fork of [masterzen/winrm](https://github.com/masterzen/winrm).<br>
If the original library gets more maintenance, I will think about switching back.

## Inspirations
* [hashicorp - terraform-provider-ad](https://github.com/hashicorp/terraform-provider-ad):<br>
Hashicorp made a great start with the terraform-provider-ad. Currently, it seems that the provider is not actively maintained.<br>
Beyond that, my goal was to split that provider into a library and a provider and extend its functionality with non Active-Directory systems.

## License
This project is licensed under the [Mozilla Public License Version 2.0](LICENSE).

<!-- Badges -->
[godoc badge]: https://pkg.go.dev/badge/github.com/d-strobel/gowindows
[godoc page]: https://pkg.go.dev/github.com/d-strobel/gowindows

[goreport badge]: https://goreportcard.com/badge/github.com/d-strobel/gowindows
[goreport page]: https://goreportcard.com/report/github.com/d-strobel/gowindows

[build badge]: https://github.com/d-strobel/gowindows/actions/workflows/build.yml/badge.svg
[build page]: https://github.com/d-strobel/gowindows/actions/workflows/build.yml

[convention badge]: https://img.shields.io/badge/Conventional%20Commits-1.0.0-%23FE5196?logo=conventionalcommits&logoColor=white
[convention page]: https://conventionalcommits.org
