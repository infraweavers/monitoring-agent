<h1 align="center">Monitoring Agent</h1>
<h3 align="center">aka "MAMA" the MAintainable Monitoring Agent</h3>

[![Test, Build and Release](https://github.com/infraweavers/monitoring-agent/actions/workflows/on-push.yml/badge.svg)](https://github.com/infraweavers/monitoring-agent/actions/workflows/on-push.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/infraweavers/mama)](https://goreportcard.com/report/github.com/infraweavers/mama)
[![License: MIT](https://img.shields.io/github/license/infraweavers/mama)](https://mit-license.org/)

### About

A simple, modern, maintainable and flexible monitoring agent that works cross platform.

### Status

In active development.

### Current Features

* Cross platform
* Executes a passed in script
* Continuous Integration/Delivery
* Windows MSI / Service
* Optional enforcement of script signing
* Optional enforcement of client TLS

### Future Planned Features

* `systemd` service
* Packaging for debian/ubuntu

### Simple Usage Examples

Linux:
```
$ curl -k -H "Content-Type: application/json" --data '{ "path": "/usr/bin/bash", "args": [ "-c", "echo \"Hello, World\"" ] }' https://test:secret@127.0.0.1:9000/v1/runscript
{"exitcode":0,"output":"Hello, World\n"}
```

Windows:
```
curl -k -H "Content-Type: application/json" --data "{ ""path"": ""C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe"", ""args"": [""-Command"", ""Write-Host 'Hello, World'""]  }" https://test:secret@127.0.0.1:9000/v1/runscript
{"exitcode":0,"output":"Hello, World\n"}
```

### Further Information

Including how to build and run Monitoring Agent, see the [wiki](https://github.com/infraweavers/monitoring-agent/wiki#building).

