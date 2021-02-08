install-deps:
	cd src/ && go mod download

# Standard go test
test:
	cd src/ && go test ./... -v -race

# Make sure no unnecessary dependencies are present
go-mod-tidy:
	cd src/ && go mod tidy -v
	git diff-index --quiet HEAD

format:
	cd src/ && go fmt $(go list ./... | grep -v /vendor/)
	cd src/ && go vet $(go list ./... | grep -v /vendor/)

generateChangelog:
	./tools/git-chglog_linux_amd64 --config tools/chglog/config.yml 0.0.0.. > CHANGELOG.md

build/snapshot:
	./tools/goreleaser_linux_amd64 --snapshot --rm-dist --skip-publish

build/release:
	./tools/goreleaser_linux_amd64 --rm-dist --skip-publish

build/docker:
	cd src/ && go build -o exporter \
     -ldflags '-X main.version=$(shell git describe --tags) -X main.commit=${shell git rev-parse HEAD} -X "main.date=${shell date --rfc-3339=seconds}" -X main.builtBy=docker' exporter.go

docker/build-latest:
	docker build -t registry.gitlab.com/hectorjsmith/fail2ban-prometheus-exporter:latest .

docker/build-tag:
	docker build -t registry.gitlab.com/hectorjsmith/fail2ban-prometheus-exporter:$(shell git describe --tags) .
