.PHONY: download test fmt check/dependencies check/fmt build build/docker

download:
	go mod download

test: download
	go test ./... -v -race

fmt: download
	go fmt ./...

check/dependencies: download
	go mod tidy -v
	git diff-index --quiet HEAD

check/fmt: download
	test -z $(shell gofmt -l .)

build:
	go build \
	-ldflags "\
	-X main.version=${shell git describe --tags} \
	-X main.commit=${shell git rev-parse HEAD} \
	-X main.date=${shell date --iso-8601=seconds} \
	-X main.builtBy=manual \
	" \
	-o fail2ban_exporter \
	exporter.go

build/docker: build
	docker build -t fail2ban-prometheus-exporter .
