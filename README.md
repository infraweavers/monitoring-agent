<h1 align="center">Monitoring Agent</h1>

[![Test, Build and Release](https://github.com/infraweavers/monitoring-agent/actions/workflows/on-push.yml/badge.svg)](https://github.com/infraweavers/monitoring-agent/actions/workflows/on-push.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/infraweavers/monitoring-agent)](https://goreportcard.com/report/github.com/infraweavers/monitoring-agent)
[![License: MIT](https://img.shields.io/github/license/infraweavers/monitoring-agent)](https://mit-license.org/)

A simple, modern, maintainable and flexible monitoring agent that works cross platform.

### Status

Released

### Current Features

* Designed to work with Nagios, Naemon and other similar monitoring platforms that perform active checks
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

#### Linux against Linux
```
$ curl -k -H "Content-Type: application/json" --data '{ "path": "perl", "args": [ "-e", "print \"Hello, World\"" ] }' https://test:secret@127.0.0.1:9000/v1/runexecutable
```

#### Windows (cmd) against Linux
```
curl -k -H "Content-Type: application/json" --data "{""path"":""perl"",""args"":[""-e"",""print 'Hello, World'""]}" https://test:secret@127.0.0.1:9000/v1/runexecutable
```

#### Windows (cmd) against Windows
```
curl -k -H "Content-Type: application/json" --data "{""path"":""C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe"",""args"":[""-Command"",""Write-Host Hello, World""]}" https://test:secret@127.0.0.1:9000/v1/runexecutable
```

#### Windows (powershell 7) against Windows
```
Invoke-RestMethod -SkipCertificateCheck -Method POST -UseBasicParsing -Credential (Get-Credential) -Uri "https://127.0.0.1:9000/v1/runexecutable" -ContentType "application/json" -Body '{"path":"C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe","args":["-Command","Write-Host Hello, World"]}'
```

#### Linux against Windows
```
curl -k -H "Content-Type: application/json" --data '{ "path": "C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe", "args": [ "-Command", "Write-Host Hello, World" ] }' https://test:secret@127.0.0.1:9000/v1/runexecutable
```
All of these examples should produce the following output.
Output:
```
{"exitcode":0,"output":"Hello, World\n"}
```

### Integration with Nagios/Naemon etc.

See the [Monitoring Agent Scripts Repository](https://github.com/infraweavers/monitoring-agent-scripts)

For the performance optimised `check_nrpe` equivalent see the [Monitoring Agent Client Repository](https://github.com/infraweavers/monitoring-agent-client)

### Wiki

For more information, including how to build and run Monitoring Agent, see the [wiki](https://github.com/infraweavers/monitoring-agent/wiki#building).

**Read the [security page](https://github.com/infraweavers/monitoring-agent/wiki/Security)** *before* using Monitoring Agent in a production environment. The default configuration has been created for ease of testing and is inherently insecure.
