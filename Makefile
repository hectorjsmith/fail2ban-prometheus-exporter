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
	mkdir -p build/
	go build \
	-ldflags "\
	-X main.version=${shell git describe --tags} \
	-X main.commit=${shell git rev-parse HEAD} \
	-X main.date=${shell date --iso-8601=seconds} \
	-X main.builtBy=manual \
	" \
	-o build/fail2ban_exporter \
	exporter.go
