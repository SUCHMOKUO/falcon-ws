#!/usr/bin/env bash

unameOut="$(uname -s)"
case "${unameOut}" in
    Linux*)     machine=linux;;
    Darwin*)    machine=mac;;
esac

CGO_ENABLED=0 go build -o "./bin/client-"$machine ./falcon-ws-client/main.go

go build -o "./bin/server-"$machine -ldflags '-linkmode "external" -extldflags "-static"' ./falcon-ws-server/main.go