#!/usr/bin/env bash

# script for Linux (WSL-compatible) dev workflow
# @Brandon Blanker Lim-it
# Thanks to @eihigh at Ebitengine's Discord server

set -euf -o pipefail

deps() {
	go get -tool github.com/dkorunic/betteralign/cmd/betteralign@latest
	go get -tool go.uber.org/nilaway/cmd/nilaway@latest
	go get -tool github.com/alexkohler/prealloc@latest
	go get -tool github.com/nikolaydubina/smrcptr@latest
	go get -tool github.com/mdempsky/unconvert@latest
	go get -tool github.com/kisielk/errcheck@latest
}

sc() {
	local -; set -x;
	go fmt ./...
	go mod tidy
	go vet ./...

	go tool betteralign -apply ./...
	go tool nilaway ./...
	go tool prealloc ./...
	go tool smrcptr ./...
	go tool unconvert ./...
	go run golang.org/x/tools/gopls/internal/analysis/modernize/cmd/modernize@latest -test ./...

	set +f
	local gofiles=( cmd/*.go src/**/*.go )
	for file in "${gofiles[@]}"; do
		goimports -w -local -v "$file"
	done
	set -f

	local PKGS
	PKGS=$(go list ./... | tr "\n" " ")
	go tool errcheck $PKGS

	set +f
	local GODIRS
	GODIRS=$(go list -f "{{.Dir}}" ./...)
	for d in ${GODIRS}; do
		go tool goimports -w -local -v "$d"/*.go
	done
	set -f
}

runwin() {
	# GOOS=windows \
	# 	go build -o out/game.exe ./cmd/main.go \
	# 	&& cp out/game.exe /mnt/c/Users/flame/game.exe \
	# 	&& /mnt/c/Users/flame/game.exe --dev
	GOOS=windows go run ./cmd/main.go --dev
}

runlinux() {
	GOOS=linux go run ./cmd/main.go --dev
}

run() {
	local -; set -x;
	if ! grep -qi Microsoft /proc/version; then
		runlinux
	else
		runwin
	fi
}

if [ "$#" -eq 0 ]; then
	echo "First use: chmod +x run.sh"
	echo "Usage: ./run.sh run"
else
	echo "Running ${1}..."
	time "$1" "$@"
fi
