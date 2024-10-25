#!/usr/bin/env bash

# script for Linux (WSL-compatible) dev workflow
# @Brandon Blanker Lim-it
# Thanks to @eihigh at Ebitengine's Discord server

set -euf -o pipefail

deps() {
	go install github.com/alexkohler/prealloc@latest
	go install github.com/nikolaydubina/smrcptr@latest
	go install go.uber.org/nilaway/cmd/nilaway@latest
	go install github.com/kisielk/errcheck@latest
}

check() {
	go mod tidy
	go vet ./...
	prealloc ./...
	smrcptr ./...
	nilaway ./...
	errcheck ./...
}

fmt() {
	set +f
	local gofiles=( src/**/*.go )
	for file in "${gofiles[@]}"; do
		goimports -w -local -v "$file"
	done
	set -f
}

run() {
	GOOS=windows \
		go build -o out/game.exe ./cmd/main.go \
		&& cp out/game.exe /mnt/c/Users/flame/game.exe \
		&& /mnt/c/Users/flame/game.exe --dev
}

if [ "$#" -eq 0 ]; then
	echo "First use: chmod +x run.sh"
	echo "Usage: ./run.sh run"
else
	echo "Running $@"
	"$1" "$@"
fi
