#!/bin/bash
echo "Building viki for ARM..."
GOOS=linux GOARCH=arm GOARM=6 CGO_ENABLED=0 go build ../vikid

scp vikid chip@chip:~/viki/
scp objects.conf chip@chip:~/viki/
scp ../resources/* chip@chip:~/viki/resources/

ssh chip@chip "kill `pgrep viki`"

