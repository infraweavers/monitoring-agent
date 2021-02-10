<h1 align="center">MAMA</h1>
<h3 align="center">MAintainable Monitoring Agent</h3>

[![Build Status](https://github.com/infraweavers/mama/workflows/Ubuntu-Test/badge.svg)](https://github.com/infraweavers/mama/actions?query=workflow:Ubuntu-Test)
[![Build Status](https://github.com/infraweavers/mama/workflows/Windows-Test/badge.svg)](https://github.com/infraweavers/mama/actions?query=workflow:Windows-Test)

### About

MAMA aims to be a flexible, easy to maintain monitoring agent that works cross platform. It is currently in active development and an official release candidate should be available soon.

### Example of current usage:

```
$ curl -k -H "Content-Type: application/json" --data '{ "path": "C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe", "args": [ "-command", "write-host \"Hello, World\"" ] }' https://test:secret@10.2.16.142:9000/v1/runscript
{"exitcode":0,"output":"Hello, World\n"}

```

### Status

In active development.

### Current Features

* Execute a passed in script

### Planned Features (active development)

* Continuous Integration/Delivery
* Tests

### Future Planned Features

* Configuration mechanism
* Windows service
* Windows MSI
* Script signing enforcement
* `systemd` service
* Packaging for debian/ubuntu
* Cross platform

### Compiling and Executing

#### Windows

1. Install [GO](https://golang.org/doc/install)
1. `git clone `https://github.com/infraweavers/mama %GOPATH%\src\mama`
1. `go run %GOPATH%/src/mama/cmd/mamasrv` 
