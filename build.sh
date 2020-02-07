#!/usr/bin/env sh

type=${1:-"client"}
target=${2:-"$(go env var GOOS | tr -d '\n')"}
mode=${3:-"development"}

output="./bin/$type-$target"
if [ "$target" = "windows" ]; then
  output="$output.exe"
fi

command="CGO_ENABLED=0 GOOS=$target go build -o $output"

if [ "$type" = "server" ]; then
  command="$command -tags='server'"
fi

if [ "$mode" = "production" ]; then
  command="$command -ldflags='-w -s'"
fi

command="$command ./falcon-ws-$type/main.go"

echo "runing: $command"
sh -c "$command"