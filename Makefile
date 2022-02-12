go/dependencies:
	cd src/ && go mod download

# Make sure no unnecessary dependencies are present
go/checkDependencies:
	cd src/ && go mod tidy -v
	git diff-index --quiet HEAD

# Standard go test
go/test:
	cd src/ && go test ./... -v -race

go/fmt:
	cd src/ && go fmt ./...

go/checkFmt:
	test -z $(shell gofmt -l .)

docs/genChangelog:
	./tools/git-chglog_linux_amd64 --config tools/chglog/config.yml 0.0.0.. > CHANGELOG_gen.md

build/snapshot:
	./tools/goreleaser_linux_amd64 --snapshot --rm-dist --skip-publish

build/release:
	./tools/goreleaser_linux_amd64 --rm-dist --skip-publish

build/docker:
	cd src/ && go build -o exporter \
     -ldflags '-X main.version=$(shell git describe --tags) -X main.commit=${shell git rev-parse HEAD} -X "main.date=${shell date --rfc-3339=seconds}" -X main.builtBy=docker' exporter.go

docker/build/latest:
	docker build -t registry.gitlab.com/hectorjsmith/fail2ban-prometheus-exporter:latest .

docker/build/nightly:
	docker build -t registry.gitlab.com/hectorjsmith/fail2ban-prometheus-exporter:nightly .

docker/build/tag:
	docker build -t registry.gitlab.com/hectorjsmith/fail2ban-prometheus-exporter:$(shell git describe --tags) .
