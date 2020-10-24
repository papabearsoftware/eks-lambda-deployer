#!/usr/bin/env bash

GOOS=linux GOARCH=amd64 CGOENABLED=0 go build -o app ./main.go

docker build -t $1:$(date +%s) .