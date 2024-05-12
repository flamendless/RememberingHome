#!/usr/bin/env bash

# script for Linux (WSL-compatible) dev workflow
# @Brandon Blanker Lim-it
# Thanks to @eihigh at Ebitengine's Discord server

set -euf -o pipefail

run() {
	GOOS=windows \
		go build -o out/game.exe ./cmd/main.go \
		&& cp out/game.exe /mnt/c/Users/flame/game.exe \
		&& /mnt/c/Users/flame/game.exe --dev
}

if [ "$#" -eq 0 ]; then
	echo "First use: chmod +x build.sh"
	echo "Usage: ./build.sh run"
else
	"$1" "$@"
fi
