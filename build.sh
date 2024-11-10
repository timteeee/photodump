#!/usr/bin/env bash

bunx tailwindcss -i css/input.css -o public/css/styles.css
go "$@"
