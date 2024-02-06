#!/usr/bin/env bash

GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc go build -o out/agg-windows-amd64.exe
GOOS=linux GOARCH=amd64 CGO_ENABLED=1 CC=gcc go build -o out/agg-linux-amd64