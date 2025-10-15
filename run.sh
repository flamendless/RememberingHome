#!/usr/bin/env bash

# script for Linux (WSL-compatible) dev workflow
# @Brandon Blanker Lim-it
# Thanks to @eihigh at Ebitengine's Discord server

set -euf -o pipefail

ISMAC=false
ISWSL=false

if [[ $(uname) == "Darwin" ]]; then
	ISMAC=true
elif grep -qi Microsoft /proc/version; then
	ISWSL=true
fi

sc() {
	local -; set -x;

	go fmt ./...
	go mod tidy
	go vet ./...

	if [ -x "$(command -v golangci-lint)" ]; then
		golangci-lint config verify
		golangci-lint run
	fi

	go tool betteralign -apply ./...
	go tool nilaway ./...
	go tool smrcptr ./...
	go tool unconvert ./...
	go run golang.org/x/tools/gopls/internal/analysis/modernize/cmd/modernize@latest -test ./...

	set +f
	local gofiles
	gofiles=( cmd/*.go src/**/*.go )
	for file in "${gofiles[@]}"; do
		go tool goimports -w -local -v "$file"
	done
	set -f

	set +f
	local GODIRS
	GODIRS=$(go list -f "{{.Dir}}" ./...)
	for d in ${GODIRS}; do
		go tool goimports -w -local -v "$d"/*.go
	done
	set -f
}

runwin() {
	GOOS=windows go run ./cmd/main.go --dev
}

runlinux() {
	GOOS=linux go run ./cmd/main.go --dev
}

runmac() {
	go run ./cmd/main.go --dev
}

runwasm() {
	local pid
	pid=$(lsof -ti:8080 2>/dev/null || true)

	if [ -n "$pid" ]; then
		local proc_name
		proc_name=$(ps -p "$pid" -o command= 2>/dev/null || true)

		if [[ "$proc_name" == *"wasmserve"* ]]; then
			echo "Killing existing wasmserve process (PID: $pid)..."
			kill "$pid"
			sleep 1
		fi
	fi

	BROWSER="${BROWSER:-vivaldi}"
	if "${ISWSL}"; then
		cmd.exe /c "start ${BROWSER} http://localhost:8080"
	elif "${ISMAC}"; then
		open -a ${BROWSER} "http://localhost:8080"
	fi

	go run github.com/hajimehoshi/wasmserve@latest ./cmd/main.go --dev
}

run() {
	local -; set -x;
	if "${ISWSL}"; then
		runwin
	elif "${ISMAC}"; then
		runmac
	else
		runlinux
	fi
}

build() {
	local -; set -x;
	go build -o "main" ./cmd/main.go
}

if [ "$#" -eq 0 ]; then
	echo "First use: chmod +x run.sh"
	echo "Usage: ./run.sh run"
else
	echo "Running ${1}..."
	time "$1" "$@"
fi
