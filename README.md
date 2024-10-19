# gURL

## Overview

This application `gURL` (you can read it girl for fun) is an attempt to mimic `cURL` app. It's written in Go but without the help of `net/http` package.

This is purely for educational purposes to understand the fundamentals of networking, web security, and HTTP.

## Set up

You must make user that there is an environment variable for the path to ceritificates file (`.pem`). This is necessary for HTTPS connections.

```bash
export CERTS_DIR='your/path/to/certs_file.pem'
```

To obtain the certificates file, you can use `homebrew`.

```bash
brew install ca-certificates
```

## Usage:

```bash
go run cmd/main.go YOUR_DOMAIN [-flags]
```

## Example:

```bash
go run cmd/main.go google.com

go run cmd/main.go https://swapi.dev/api/people/1 -json
```
