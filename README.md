# Fail2Ban Prometheus Exporter

[![Pipeline](https://gitlab.com/hectorjsmith/fail2ban-prometheus-exporter/badges/main/pipeline.svg)](https://gitlab.com/hectorjsmith/fail2ban-prometheus-exporter)

Collect metrics from a running fail2ban instance.

## Table of Contents
1. Quick Start
2. Metrics
3. Configuration
4. Building from source
5. Textfile metrics

## 1. Quick Start

The exporter can be run as a standalone binary or a docker container.

### 1.1. Standalone

The following command will start collecting metrics from the `/var/run/fail2ban/fail2ban.sock` file and expose them on port `9191`.

```
$ fail2ban_exporter --collector.f2b.socket=/var/run/fail2ban/fail2ban.sock --web.listen-address=":9191"

2022/02/20 09:54:06 fail2ban exporter version 0.5.0
2022/02/20 09:54:06 starting server at :9191
2022/02/20 09:54:06 reading metrics from fail2ban socket: /var/run/fail2ban/fail2ban.sock
2022/02/20 09:54:06 metrics available at '/metrics'
2022/02/20 09:54:06 ready
```

Binary files for each release can be found on the [releases](https://gitlab.com/hectorjsmith/fail2ban-prometheus-exporter/-/releases) page.

### 1.2. Docker

**Docker run**
```
docker run -d \
    --name "fail2ban-exporter" \
    -v /var/run/fail2ban:/var/run/fail2ban:ro \
    -p "9191:9191" \
    registry.gitlab.com/hectorjsmith/fail2ban-prometheus-exporter:latest
```

**Docker compose**

```
version: "2"
services:
  exporter:
    image: registry.gitlab.com/hectorjsmith/fail2ban-prometheus-exporter:latest
    volumes:
    - /var/run/fail2ban/:/var/run/fail2ban:ro
    ports:
    - "9191:9191"
```

Use the `:latest` tag to get the latest stable release. Or use the `:nightly` tag for the latest (unstable) version.
See the [registry page](https://gitlab.com/hectorjsmith/fail2ban-prometheus-exporter/container_registry) for all available tags.

**NOTE:** While it is possible to mount the `fail2ban.sock` file directly, it is recommended to mount the parent folder instead.
The `.sock` file is deleted by fail2ban on shutdown and re-created on startup and this causes problems for the docker mount.
See [this reply](https://gitlab.com/hectorjsmith/fail2ban-prometheus-exporter/-/issues/11#note_665003499) for more details.

## 2. Metrics

The exporter exposes the following metrics:

*All metric names are prefixed with `f2b_`*

| Metric                       | Description                                                                        | Example                                             |
|------------------------------|------------------------------------------------------------------------------------|-----------------------------------------------------|
| `up`                         | Returns 1 if the exporter is up and running                                        | `f2b_up 1`                                          |
| `errors`                     | Count the number of errors since startup by type                                   |                                                     |
| `errors{type="socket_conn"}` | Errors connecting to the fail2ban socket (e.g. connection refused)                 | `f2b_errors{type="socket_conn"} 0`                  |
| `errors{type="socket_req"}`  | Errors sending requests to the fail2ban server (e.g. invalid responses)            | `f2b_errors{type="socket_req"} 0`                   |
| `jail_count`                 | Number of jails configured in fail2ban                                             | `f2b_jail_count 2`                                  |
| `jail_banned_current`        | Number of IPs currently banned per jail                                            | `f2b_jail_banned_current{jail="sshd"} 15`           |
| `jail_banned_total`          | Total number of banned IPs since fail2ban startup per jail (includes expired bans) | `f2b_jail_banned_total{jail="sshd"} 31`             |
| `jail_failed_current`        | Number of current failures per jail                                                | `f2b_jail_failed_current{jail="sshd"} 6`            |
| `jail_failed_total`          | Total number of failures since fail2ban startup per jail                           | `f2b_jail_failed_total{jail="sshd"} 125`            |
| `jail_config_ban_time`       | How long an IP is banned for in this jail (in seconds)                             | `f2b_config_jail_ban_time{jail="sshd"} 600`         |
| `jail_config_find_time`      | How far back the filter will look for failures in this jail (in seconds)           | `f2b_config_jail_find_time{jail="sshd"} 600`        |
| `jail_config_max_retry`      | The max number of failures allowed before banning an IP in this jail               | `f2b_config_jail_max_retries{jail="sshd"} 5`        |
| `version`                    | Version string of the exporter and fail2ban                                        | `f2b_version{exporter="0.5.0",fail2ban="0.11.1"} 1` |

The metrics above correspond to the matching fields in the `fail2ban-client status <jail>` command:
```
Status for the jail: sshd
|- Filter
|  |- Currently failed: 6
|  |- Total failed:     125
|  `- File list:        /var/log/auth.log
`- Actions
   |- Currently banned: 15
   |- Total banned:     31
   `- Banned IP list:   ...
```

### 2.1. Grafana

The metrics exported by this tool are compatible with Prometheus and Grafana.
A sample grafana dashboard can be found in the [grafana.json](/examples/grafana/dashboard.json) file.
Just import the contents of this file into a new Grafana dashboard to get started.

*(Sample dashboard is compatible with Grafana `8.3.3` and above)*

## 3. Configuration

The exporter is configured with CLI flags and environment variables.
There are no configuration files.

**CLI flags**
```
usage: exporter [<flags>]

Flags:
  -h, --help     Show context-sensitive help (also try --help-long and --help-man).
  -v, --version  show version info and exit
      --collector.f2b.socket="/var/run/fail2ban/fail2ban.sock"  
                 path to the fail2ban server socket
      --collector.textfile.directory=""  
                 directory to read text files with metrics from
      --web.listen-address=":9191"  
                 address to use for the metrics server
      --web.basic-auth.username=""  
                 username to use to protect endpoints with basic auth
      --web.basic-auth.password=""  
                 password to use to protect endpoints with basic auth
      --collector.f2b.exit-on-socket-connection-error  
                 when set to true the exporter will immediately exit on a fail2ban socket connection error
```

**Environment variables**

Each environment variable corresponds to a CLI flag.
If both are specified, the CLI flag takes precedence.

| Environment variable       | Corresponding CLI flag                            |
|----------------------------|---------------------------------------------------|
| `F2B_COLLECTOR_SOCKET`     | `--collector.f2b.socket`                          |
| `F2B_COLLECTOR_TEXT_PATH`  | `--collector.textfile.directory`                  |
| `F2B_WEB_LISTEN_ADDRESS`   | `--web.listen-address`                            |
| `F2B_WEB_BASICAUTH_USER`   | `--web.basic-auth.username`                       |
| `F2B_WEB_BASICAUTH_PASS`   | `--web.basic-auth.password`                       |
| `F2B_EXIT_ON_SOCKET_ERROR` | `--collector.f2b.exit-on-socket-connection-error` |

## 4. Building from source

The simplest way to build the project is to run the `build/snapshot` make command.
This will use `goreleaser` to build out binaries and archives for the project.
Binaries are stored in the `dist/` folder.

Alternatively, `go mod download` and `go build` can be used from the `src/` folder to build out the project.
This will download dependencies and build the project.

## 5. Textfile metrics

For more flexibility the exporter also allows exporting metrics collected from a text file.

To enable textfile metrics provide the directory to read files from with the `--collector.textfile.directory` flag.

Metrics collected from these files will be exposed directly alongside the other metrics without any additional processing.
This means that it is the responsibility of the file creator to ensure the format is correct.

By exporting textfile metrics an extra metric is also exported with an error count for each file:

```
# HELP textfile_error Checks for errors while reading text files
# TYPE textfile_error gauge
textfile_error{path="file.prom"} 0
```

**NOTE:** Any file not ending with `.prom` will be ignored.

**Running in Docker**

To collect textfile metrics inside a docker container, a couple of things need to be done:
1. Mount the folder with the metrics files
2. Set the `F2B_COLLECTOR_TEXT_PATH` environment variable

*For example:*
```
docker run -d \
    --name "fail2ban-exporter" \
    -v /var/run/fail2ban:/var/run/fail2ban:ro \
    -v /path/to/metrics:/app/metrics/:ro \
    -e F2B_COLLECTOR_TEXT_PATH=/app/metrics \
    -p "9191:9191" \
    registry.gitlab.com/hectorjsmith/fail2ban-prometheus-exporter:latest
```
