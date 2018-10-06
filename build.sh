#!/bin/bash

export GOOS=$1
export GOARCH=$2
export GOPATH=`pwd`:$GOPATH

go install falcon-ws-server
go install falcon-ws-client
