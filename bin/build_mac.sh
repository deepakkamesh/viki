#!/bin/bash
echo "Building viki for mac..."
GOOS=darwin GOARCH=amd64 go build ../vikid
./vikid
