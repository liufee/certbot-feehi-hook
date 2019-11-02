#!/bin/bash

CGO_ENABLED=0 go build -o certbot-feehi-hook -a -ldflags '-extldflags "-static"' .