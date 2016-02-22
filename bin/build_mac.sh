#!/bin/bash
BUILDTIME="`date '+%Y-%m-%d_%I:%M:%S%p'`"
GITHASH="`git rev-parse --short=7 HEAD`"
VER="-X main.buildtime=$BUILDTIME -X main.githash=$GITHASH"
echo "Building viki version:$GITHASH for mac"
GOOS=darwin GOARCH=amd64 go build -ldflags "$VER" ../vikid
