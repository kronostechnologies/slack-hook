#!/bin/bash

CGO_ENABLED=0 go build -ldflags="-w -s" -o slack-hook .
upx --best slack-hook
