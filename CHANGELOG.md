# CHANGELOG

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog], and this project adheres to [Semantic Versioning].

## [Unreleased]

### Added
- (23e073f) feat: add example systemd service file
- (3911eca) feat: rename output binary and archives
- (f6e328a) feat: correctly handle shutdown signals
- (6e575aa) feat: rewrite cli flags and environment variables
- (0f0efe5) feat: remove startup script from docker image
- (e2902b8) feat: improve logging on startup
- (1f27dac) feat: add default value for fail2ban socket path
- (b7e317e) feat: configure tool using environment variables ([#17](https://gitlab.com/hectorjsmith/fail2ban-prometheus-exporter/issues/17))
- (6f76a03) feat: add support for basic auth ([#16](https://gitlab.com/hectorjsmith/fail2ban-prometheus-exporter/issues/16))

### Fixed
- (93da909) fix: use correct flag in dockerfile ([#18](https://gitlab.com/hectorjsmith/fail2ban-prometheus-exporter/issues/18))

### BREAKING CHANGE
- Release binary name has been changed to `fail2ban_exporter`.
- Replace `--socket` flag with `--collector.f2b.socket`.
- Merge `--port` flag and `--web.listen-address` into a single flag.
- Remove `--collector.textfile` flag, its value is now derived from `--collector.textfile.directory`.
- Remove `F2B_COLLECTOR_TEXT` and `F2B_WEB_PORT` environment variables.
- Using the textfile collector in docker now requires setting environment variables.
- CLI params now require two dashes instead of one (e.g. `--socket`).

## [0.5.0] - 2021-12-21
*Remove deprecated code & support python2*

### Added
- (b397a51) feat: sample grafana dashboard ([#15](https://gitlab.com/hectorjsmith/fail2ban-prometheus-exporter/issues/15))
- (c208c8e) feat: add listen address parameter - thanks [@private-creator](https://gitlab.com/private-creator)!

### Fixed
- (7932ccb) fix: support python2 fail2ban ([#14](https://gitlab.com/hectorjsmith/fail2ban-prometheus-exporter/issues/14)) - thanks [@private-creator](https://gitlab.com/private-creator)!

### BREAKING CHANGE
- Remove `-db` CLI flag
- Remove `f2b_errors{type="db"}` metric

## [0.4.0] - 2021-10-18
*Add new fail2ban config metrics*

### Added
- (56730c8) feat: add new jail config metrics
- (5a107cc) feat: support for textfile metrics ([#13](https://gitlab.com/hectorjsmith/fail2ban-prometheus-exporter/issues/13))

### Removed
- (b268f86) remove: database-based metrics
- (0b6a941) remove: windows builds

## [0.3.0] - 2021-09-27
*Export new version metrics ([#12](https://gitlab.com/hectorjsmith/fail2ban-prometheus-exporter/issues/12))*

### Added
- (3c9a005) feat: render basic html page at root url
- (22a165d) feat: improve startup logging
- (fba9ee2) feat: export new version metric ([#12](https://gitlab.com/hectorjsmith/fail2ban-prometheus-exporter/issues/12))

## [0.2.0] - 2021-08-31
*Collect metrics through fail2ban socket - based on [#11](https://gitlab.com/hectorjsmith/fail2ban-prometheus-exporter/issues/11)*

### Added
- (39133d0) feat: collect new up metric from fail2ban socket
- (4da46f3) feat: export metrics with socket errors
- (bd841c3) feat: set up metric to 0 if errors found
- (1964dde) feat: export metrics for failed/banned counts
- (2ab1f7d) feat: support reading fail2ban socket in docker
- (1282d63) feat: new metric for enabled jails ([#1](https://gitlab.com/hectorjsmith/fail2ban-prometheus-exporter/issues/1))

### Fixed
- (526b1c7) fix: update banned metrics to exclude expired bans ([#11](https://gitlab.com/hectorjsmith/fail2ban-prometheus-exporter/issues/11))

### Deprecated
- Use of the fail2ban database has been deprecated. The exporter now collects metrics through the fail2ban socket file. See [#11](https://gitlab.com/hectorjsmith/fail2ban-prometheus-exporter/-/issues/11) for more details.

## [0.1.0] - 2021-03-28
*Initial release*

### Added
- (6355c9e) feat: fail on startup if database file does not exist ([#8](https://gitlab.com/hectorjsmith/fail2ban-prometheus-exporter/issues/8))
- (4f18bf3) feat: add cli parameters for db path and metrics port ([#4](https://gitlab.com/hectorjsmith/fail2ban-prometheus-exporter/issues/4))
- (91cba80) feat: export number of banned ips
- (4b96501) feat: export bad ip count per jail
- (0b40e5d) feat: connect to fail2ban db and extract total bad ips
- (7ced846) feat: initial setup of metric exporter


### Fixed
- (0842419) fix: compile tool without cgo_enabled flag

## 0.0.0 - 2021-02-05
*Repository creation*

---

[Keep a Changelog]: https://keepachangelog.com/en/1.0.0/
[Semantic Versioning]: https://semver.org/spec/v2.0.0.html
[Unreleased]: https://gitlab.com/hectorjsmith/fail2ban-prometheus-exporter/compare/0.1.0...main
[0.1.0]: https://gitlab.com/hectorjsmith/fail2ban-prometheus-exporter/compare/0.0.0...0.1.0
[0.2.0]: https://gitlab.com/hectorjsmith/fail2ban-prometheus-exporter/compare/0.1.0...0.2.0
[0.3.0]: https://gitlab.com/hectorjsmith/fail2ban-prometheus-exporter/compare/0.2.0...0.3.0
[0.4.0]: https://gitlab.com/hectorjsmith/fail2ban-prometheus-exporter/compare/0.3.0...0.4.0
[0.5.0]: https://gitlab.com/hectorjsmith/fail2ban-prometheus-exporter/compare/0.4.0...0.5.0
