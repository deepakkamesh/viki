#!/bin/bash
echo "Building viki for ARM"
GOOS=linux GOARCH=arm GOARM=6 CGO_ENABLED=0 go build ../main
scp main pi@10.0.0.23:~/
