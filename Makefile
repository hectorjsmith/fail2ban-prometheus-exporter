# List make commands
.PHONY: ls
ls:
	cat Makefile | grep "^[a-zA-Z#].*" | cut -d ":" -f 1 | sed s';#;\n#;'g

# Download dependencies
.PHONY: download
download:
	go mod download

# Update project dependencies
.PHONY: update
update:
	go get -u
	go mod download
	go mod tidy

# Run project tests
.PHONY: test
test: download
	go test ./... -v -race

# Format code
.PHONY: fmt
fmt: download
	go mod tidy
	go fmt ./...

# Check for unformatted go code
.PHONY: check/fmt
check/fmt: download
	test -z $(shell gofmt -l .)

# Build project
.PHONY: build
build:
	CGO_ENABLED=0 go build \
	-ldflags "\
	-X main.version=${shell git describe --tags} \
	-X main.commit=${shell git rev-parse HEAD} \
	-X main.date=${shell date --iso-8601=seconds} \
	-X main.builtBy=manual \
	" \
	-o fail2ban_exporter \
	exporter.go

# Build project docker container
.PHONY: build/docker
build/docker: build
	docker build -t fail2ban-prometheus-exporter .
