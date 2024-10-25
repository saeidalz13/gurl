# gURL

## Overview

This application `gURL` (you can read it girl for fun) is an attempt to mimic `cURL` app. It's written in Go but without the help of `net/http` package.

This is purely for educational purposes to understand the fundamentals of networking, web security, and HTTP.

## Usage:

```bash
go run cmd/main.go YOUR_DOMAIN [-flags]
```

For websocket connetions, you **must** include the protocol.

```bash
go run cmd/main.go ws://YOUR_DOMAIN [-flags]
```

```bash
go run cmd/main.go wss://YOUR_DOMAIN [-flags]
```

## Example:

```bash
go run cmd/main.go google.com

go run cmd/main.go https://swapi.dev/api/people/1 -json
```
