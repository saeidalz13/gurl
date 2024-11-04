# gURL

## Overview

This application `gURL` (you can read it girl for fun) is an attempt to mimic `cURL` app. It's written in Go but without the help of `net/http` package. \
This is mostly for educational purposes to understand the fundamentals of networking, web security, and HTTP.

`gURL` also attempts to provide a prettier terminal experience compared to `cURL`.

## Usage:

```bash
Usage app.exe DOMAIN [flags]:
  -cookies string
        Add cookie to request header; e.g. -cookies='name1=value1; name2=value2'
  -json string
        Add json data to body
  -method string
        HTTP method (default "GET")
  -text string
        Add plain text to body
  -v    Verbose run
```

## WebSocket:

For websocket connetions, you **must** include the protocol.

```bash
# Excluding TLS
go run cmd/main.go ws://YOUR_DOMAIN [-flags]
```

```bash
# Including TLS
go run cmd/main.go wss://YOUR_DOMAIN [-flags]
```

## Examples:

```bash
go run cmd/main.go www.google.com

go run cmd/main.go swapi.dev/api/people/1 -v

gocmd https://jsonplaceholder.typicode.com/posts -json='{"title":"foo","body":"bar","userId":1}' -method=post -v

```
