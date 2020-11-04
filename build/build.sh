#!/usr/bin/env bash

GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -race -o main ./cmd

zip latest.zip main