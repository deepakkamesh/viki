#!/bin/bash
echo "Building viki for ARM..."
GOOS=linux GOARCH=arm GOARM=6 CGO_ENABLED=0 go build ../vikid
scp vikid pi@10.0.0.23:~/viki/