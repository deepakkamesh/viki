#!/bin/bash
echo "Building viki for ARM..."
GOOS=linux GOARCH=arm GOARM=6 CGO_ENABLED=1 go build ../vikid
