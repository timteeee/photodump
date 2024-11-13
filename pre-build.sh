#!/usr/bin/env bash

set -e

PUBLIC_DIR="./assets/public"
HTMX_VERSION="2.0.3"
ALPINEJS_VERSION="3.14.3"

# Download HTMX if it's not there already
HTMX_FILE="$PUBLIC_DIR/js/htmx-$HTMX_VERSION.min.js" 

if ! [ -f $HTMX_FILE ]; then
	HTMX_URL="https://unpkg.com/htmx.org@$HTMX_VERSION/dist/htmx.min.js"
	echo "HTMX source file not found, fetching from $HTMX_URL..."

	curl $HTMX_URL > $HTMX_FILE
else
	echo "HTMX source file already exists, skipping this step."
fi

# Download Alpine.js if it's not there already
ALPINEJS_FILE="$PUBLIC_DIR/js/alpinejs-$ALPINEJS_VERSION.min.js" 

if ! [ -f $ALPINEJS_FILE ]; then
	ALPINEJS_URL="https://cdn.jsdelivr.net/npm/alpinejs@$ALPINEJS_VERSION/dist/cdn.min.js"
	echo "Alpine.js source file not found, fetching from $ALPINEJS_URL..."

	curl $ALPINEJS_URL > $ALPINEJS_FILE
else
	echo "Alpine.js source file already exists, skipping this step."
fi

# Run Tailwind
echo "Rebuilding Tailwind styles..."
bunx tailwindcss -i ./build/tailwind-input.css -o $PUBLIC_DIR/css/styles.css
