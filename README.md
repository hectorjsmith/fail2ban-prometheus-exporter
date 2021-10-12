# Fail2Ban Prometheus Exporter

Go tool to collect and export metrics on Fail2Ban

## Table of Contents
1. Introduction
2. Running the Exporter
3. Running in Docker
4. Metrics

## 1. Introduction
The exporter can collect metrics from 2 locations: the fail2ban server socket and the fail2ban server database.

Once the exporter is running, metrics are available at `localhost:9191/metrics`.

(The default port is `9191` but can be modified with the `-port` flag)

### 1.1. Socket
The recommended way to run the exporter is to point it at the fail2ban server socket.
This allows the exporter to communicate with the server the same way `fail2ban-client` does and ensures the metrics it collects align with the values reported by `fail2ban-client status <jail>`.

The default path to the socket is: `/var/run/fail2ban/fail2ban.sock`

### 1.2. Deprecated: Database
The original way to collect metrics is to read them from the fail2ban database.
This has now been deprecated in favour of using the socket.
The reason being that database metrics do not always align with the output of `fail2ban-client status <jail>` and cause confusion.
See [#11](https://gitlab.com/hectorjsmith/fail2ban-prometheus-exporter/-/issues/11) for more details.

If necessary, these metrics can still be exported by providing the database path to the exporter.

The default path to the fail2ban database is: `/var/lib/fail2ban/fail2ban.sqlite3`

## 2. Running the Exporter

The exporter is compiled and released as a single binary.
This makes it very easy to run in any environment.
No additional runtime dependencies are required.

Compiled binaries for various platforms are provided in each release.
See the [releases page](https://gitlab.com/hectorjsmith/fail2ban-prometheus-exporter/-/releases) for more information.

**Usage**
```
$ fail2ban-prometheus-exporter -h

  -db string
        path to the fail2ban sqlite database (deprecated)
  -port int
        port to use for the metrics server (default 9191)
  -socket string
        path to the fail2ban server socket
  -version
        show version info and exit
  -collector.textfile
        enable the textfile collector
  -collector.textfile.directory string
        directory to read text files with metrics from
```

**Example**

```
fail2ban-prometheus-exporter -socket /var/run/fail2ban/fail2ban.sock -port 9191
```

Note that the exporter will need read access to the fail2ban socket or database.

### 2.1. Compile from Source

The code can be compiled from source by running `go build` inside the `src/` folder.
Go version `1.15` or greater is required.

Run `go mod download` to download all necessary dependencies before running the build.

## 3. Running in Docker

If use of docker is desired, an official docker image is available on the Gitlab container registry.
Use it by pulling the following image:

```
registry.gitlab.com/hectorjsmith/fail2ban-prometheus-exporter:latest
```

Use the `:latest` tag to get the most up to date code (less stable) or use one of the version tagged images to use a specific release.
See the [registry page](https://gitlab.com/hectorjsmith/fail2ban-prometheus-exporter/container_registry) for all available tags.

### 3.1. Volumes

The docker image is designed to run by mounting either the fail2ban sqlite3 database of the fail2ban run folder.
- The database should be mounted at: `/app/fail2ban.sqlite3`
- The run folder should be mounted at: `/var/run/fail2ban`

Both paths can be mounted with readonly (`ro`) permissions.

**NOTE:** While it is possible to mount the `fail2ban.sock` file directly, it is recommended to mount the parent folder instead.
The `.sock` file is deleted by fail2ban on shutdown and re-created on startup and this causes problems for the docker mount.
See [this reply](https://gitlab.com/hectorjsmith/fail2ban-prometheus-exporter/-/issues/11#note_665003499) for more details.

### 3.2. Docker run

Use the following command to run the exporter as a docker container.

```
docker run -d \
    --name "fail2ban-exporter" \
    -v /var/lib/fail2ban/fail2ban.sqlite3:/app/fail2ban.sqlite3:ro \
    -v /var/run/fail2ban:/var/run/fail2ban:ro \
    -p "9191:9191" \
    registry.gitlab.com/hectorjsmith/fail2ban-prometheus-exporter:latest
```

### 3.3. Docker compose

The following is a simple docker-compose file to run the exporter.

```
version: "2"
services:
  exporter:
    image: registry.gitlab.com/hectorjsmith/fail2ban-prometheus-exporter:latest
    volumes:
    - /var/lib/fail2ban/fail2ban.sqlite3:/app/fail2ban.sqlite3:ro
    - /var/run/fail2ban/:/var/run/fail2ban:ro
    ports:
    - "9191:9191"
```

## 4. Metrics

Access exported metrics at the `/metrics` path on the configured port.

**Note on Fail2Ban Jails**

fail2ban can be configured to process different log files and use different rules for each one.
These separate configurations are referred to as *jails*.

For example, fail2ban can be configured to watch the system logs for failed SSH connections and Nextcloud logs for failed logins.
In this configuration, there will be two jails - one for IPs banned from the SSH logs, and one for IPs banned from the Nextcloud logs.

This tool exports several metrics *per jail*, meaning that it is possible to track how many IPs are being banned in each jail as well as the overall total.
This can be useful to track what services are seeing more failed logins.

### 4.1. Socket-based Metrics

These are the metrics exported by reading data from the fail2ban server socket.
All metrics are prefixed with `f2b_`.

Exposed metrics:
* `up` - Returns 1 if the fail2ban server is up and connection succeeds
* `errors` - Number of errors since startup
    * `db` - Errors connecting to the database
    * `socket_conn` - Errors connecting to the fail2ban socket (e.g. connection refused)
    * `socket_req` - Errors sending requests to the fail2ban server (e.g. invalid responses)
* `jail_count` - Number of jails configured in fail2ban
* `jail_banned_current` (per jail) - Number of IPs currently banned
* `jail_banned_total` (per jail) - Total number of banned IPs since fail2ban startup (includes expired bans)
* `jail_failed_current` (per jail) - Number of current failures
* `jail_failed_total` (per jail) - Total number of failures since fail2ban startup

**Sample**

```
# HELP f2b_errors Number of errors found since startup
# TYPE f2b_errors counter
f2b_errors{type="db"} 0
f2b_errors{type="socket_conn"} 0
f2b_errors{type="socket_req"} 0
# HELP f2b_jail_banned_current Number of IPs currently banned in this jail
# TYPE f2b_jail_banned_current gauge
f2b_jail_banned_current{jail="recidive"} 5
f2b_jail_banned_current{jail="sshd"} 15
# HELP f2b_jail_banned_total Total number of IPs banned by this jail (includes expired bans)
# TYPE f2b_jail_banned_total gauge
f2b_jail_banned_total{jail="recidive"} 6
f2b_jail_banned_total{jail="sshd"} 31
# HELP f2b_jail_count Number of defined jails
# TYPE f2b_jail_count gauge
f2b_jail_count 2
# HELP f2b_jail_failed_current Number of current failures on this jail's filter
# TYPE f2b_jail_failed_current gauge
f2b_jail_failed_current{jail="recidive"} 5
f2b_jail_failed_current{jail="sshd"} 6
# HELP f2b_jail_failed_total Number of total failures on this jail's filter
# TYPE f2b_jail_failed_total gauge
f2b_jail_failed_total{jail="recidive"} 7
f2b_jail_failed_total{jail="sshd"} 125
# HELP f2b_up Check if the fail2ban server is up
# TYPE f2b_up gauge
f2b_up 1
```

The metrics above correspond to the matching fields in the `fail2ban-client status <jail>` command:
```
Status for the jail: sshd|- Filter
|  |- Currently failed: 6
|  |- Total failed:     125
|  `- File list:        /var/log/auth.log
`- Actions
   |- Currently banned: 15
   |- Total banned:     31
   `- Banned IP list:   ...
```

### 4.2. Database Metrics (deprecated)

These are the original metrics exported by the initial release of the exporter.
They are all based on the data stored in the fail2ban sqlite3 database.

*These metrics are deprecated and will be removed in a future release.*

All metrics are prefixed with `fail2ban_`.

Exposed metrics:
* `up` - Returns 1 if the service is up
* `errors` - Returns the number of errors found since startup
* `enabled_jails` - Returns 1 for each jail that is enabled, 0 if disabled.
* `bad_ips` (per jail)
    * A *bad IP* is defined as an IP that has been banned at least once in the past
    * Bad IPs are counted per jail
* `banned_ips` (per jail)
    * A *banned IP* is defined as an IP that is currently banned on the firewall
    * Banned IPs are counted per jail

**Sample**

```
# HELP fail2ban_bad_ips (Deprecated) Number of bad IPs stored in the database (per jail).
# TYPE fail2ban_bad_ips gauge
fail2ban_bad_ips{jail="recidive"} 0
fail2ban_bad_ips{jail="sshd"} 0
# HELP fail2ban_banned_ips (Deprecated) Number of banned IPs stored in the database (per jail).
# TYPE fail2ban_banned_ips gauge
fail2ban_banned_ips{jail="recidive"} 0
fail2ban_banned_ips{jail="sshd"} 0
# HELP fail2ban_enabled_jails (Deprecated) Enabled jails.
# TYPE fail2ban_enabled_jails gauge
fail2ban_enabled_jails{jail="recidive"} 1
fail2ban_enabled_jails{jail="sshd"} 1
# HELP fail2ban_errors (Deprecated) Number of errors found since startup.
# TYPE fail2ban_errors counter
fail2ban_errors{type="db"} 0
# HELP fail2ban_up (Deprecated) Was the last fail2ban query successful.
# TYPE fail2ban_up gauge
fail2ban_up 1
```

### 4.3. Textfile Metrics

For more flexibility the exporter also allows exporting metrics collected from a text file.

To enable textfile metrics:
1. Enable the collector with `-collector.textfile=true`
2. Provide the directory to read files from with the `-collector.textfile.directory` flag

Metrics collected from these files will be exposed directly alongside the other metrics without any additional processing.
This means that it is the responsibility of the file creator to ensure the format is correct.

By exporting textfile metrics an extra metric is also exported with an error count for each file:

```
# HELP textfile_error Checks for errors while reading text files
# TYPE textfile_error gauge
textfile_error{path="file.prom"} 0
```

**NOTE:** Any file not ending with `.prom` will be ignored.
