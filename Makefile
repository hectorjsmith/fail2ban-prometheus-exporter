go/dependencies:
	go mod download

# Make sure no unnecessary dependencies are present
go/checkDependencies:
	go mod tidy -v
	git diff-index --quiet HEAD

# Standard go test
go/test:
	go test ./... -v -race

go/fmt:
	go fmt ./...

go/checkFmt:
	test -z $(shell gofmt -l .)

build/docker:
	go build -o fail2ban_exporter \
     -ldflags '-X main.version=$(shell git describe --tags) -X main.commit=${shell git rev-parse HEAD} -X "main.date=${shell date --rfc-3339=seconds}" -X main.builtBy=docker' exporter.go

docker/build/latest:
	docker build -t registry.gitlab.com/hectorjsmith/fail2ban-prometheus-exporter:latest .

docker/build/nightly:
	docker build -t registry.gitlab.com/hectorjsmith/fail2ban-prometheus-exporter:nightly .

docker/build/tag:
	docker build -t registry.gitlab.com/hectorjsmith/fail2ban-prometheus-exporter:$(shell git describe --tags) .
