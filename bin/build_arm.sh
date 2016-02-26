#!/bin/bash
BUILDTIME="`date '+%Y-%m-%d_%I:%M:%S%p'`"
GITHASH="`git rev-parse --short=7 HEAD`"
# For go version < 1.5
VER="-X main.buildtime $BUILDTIME -X main.githash $GITHASH"
# For go version > 1.5
#VER="-X main.buildtime=$BUILDTIME -X main.githash=$GITHASH"
git pull
echo "Building viki version:$GITHASH for linux arm"
GOOS=linux GOARCH=arm GOARM=6 CGO_ENABLED=0 go build -ldflags "$VER" ../vikid
