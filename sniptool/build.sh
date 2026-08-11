#!/bin/sh
# Build the whole package, not a single file — snip.go was split up.
set -e
go test ./...
go build -o ../snip .
