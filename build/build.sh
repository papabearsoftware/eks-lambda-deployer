#!/usr/bin/env bash

GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -o main ./cmd

zip latest.zip main