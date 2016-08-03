#!/bin/sh

time=`date '+%Y-%m-%d-%H:%M:%S'`
version=`git log -1 --oneline`
strip_version=${version// /-}
GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w -X main.time='${time}' -X main.version='${strip_version}'" -o gister_mac
GOOS=linux  GOARCH=amd64 go build -ldflags "-s -w -X main.time='${time}' -X main.version='${strip_version}'" -o gister_linux
