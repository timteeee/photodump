#!/usr/bin/env bash

set -e

./pre-build.sh

# Pass args from gow to the Go binary
go "$@"
