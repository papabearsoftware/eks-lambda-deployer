#!/usr/bin/env bash

GOARCH=amd64 GOOS=linux CGOENABLED=0 go build -gcflags='-N -l' -o main ./cmd

zip latest.zip main