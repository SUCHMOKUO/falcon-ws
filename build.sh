#!/bin/bash

export GOOS=$1
export GOARCH=$2
export GOPATH=`pwd`:$GOPATH
export CGO_ENABLED=0

go install falcon-ws-server
go install falcon-ws-client
