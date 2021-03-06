<h1 align="center">Monitoring Agent</h1>

[![Test, Build and Release](https://github.com/infraweavers/monitoring-agent/actions/workflows/on-push.yml/badge.svg)](https://github.com/infraweavers/monitoring-agent/actions/workflows/on-push.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/infraweavers/monitoring-agent)](https://goreportcard.com/report/github.com/infraweavers/monitoring-agent)
[![License: MIT](https://img.shields.io/github/license/infraweavers/monitoring-agent)](https://mit-license.org/)

A simple, modern, maintainable and flexible monitoring agent that works cross platform.

### Status

Testing

### Current Features

* Designed to work with Nagios, Naemon and other similar monitoring platforms
* Cross platform
* Executes a passed in script (no need to deploy scripts to monitored hosts)
* Continuous Integration/Delivery
* Windows MSI / Service
* Optional enforcement of script signing
* Optional enforcement of client TLS
* `systemd` service
* Packaging for debian/ubuntu

### Features in Development

* None

### Simple Usage Examples

Linux:
```
$ curl -k -H "Content-Type: application/json" --data '{ "path": "/usr/bin/bash", "args": [ "-c", "echo \"Hello, World\"" ] }' https://test:secret@127.0.0.1:9000/v1/runscript
{"exitcode":0,"output":"Hello, World\n"}
```

Windows:
```
curl -k -H "Content-Type: application/json" --data '{ "path": "C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe", "args": ["-Command", "Write-Host Hello, World"]  }' https://test:secret@127.0.0.1:9000/v1/runscript
{"exitcode":0,"output":"Hello, World\n"}
```

### Integration with Nagios/Naemon etc.

See the [Monitoring Agent Scripts Repository](https://github.com/infraweavers/monitoring-agent-scripts)

### Wiki

For more information, including how to build and run Monitoring Agent, see the [wiki](https://github.com/infraweavers/monitoring-agent/wiki#building).

**Read the [security page](https://github.com/infraweavers/monitoring-agent/wiki/Security)** *before* using Monitoring Agent in a production environment. The default configuration has been created for ease of testing and is inherently insecure.
