#!/usr/bin/env bash

if [[ $# -eq 0 ]]
then
  echo "Repository name must be specified!"
  exit 1
fi

GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o app ./main.go

TS=$(date +%s)

docker build -t $1:$TS .

docker push $1:$TS