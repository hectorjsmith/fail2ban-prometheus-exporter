# Fail2Ban Prometheus Exporter

Go tool to collect and export metrics on Fail2Ban

## Table of Contents
1. How to use
2. Docker
    1. Environment variables
    2. Docker run
    3. Docker compose
3. CLI usage
4. Metrics

## 1. How to use

Run the exporter by providing it with a fail2ban database to read data from.
Read access to the database is required.
Once the exporter is running, metrics are available at `localhost:9191/metrics`.

The default port is `9191`, but this can be modified with the `-port` flag.

**Note:** By default fail2ban stores the database file at: `/var/lib/fail2ban/fail2ban.sqlite3`

**Fail2Ban Jails**

fail2ban can be configured to process different log files and use different rules for each one.
These separate configurations are referred to as *jails*.

For example, fail2ban can be configured to watch the system logs for failed SSH connections and Nextcloud logs for failed logins.
In this configuration, there will be two jails - one for IPs banned from the SSH logs, and one for IPs banned from the Nextcloud logs.

This tool exports several metrics *per jail*, meaning that it is possible to track how many IPs are being banned in each jail as well as the overall total.
This can be useful to track what services are seeing more failed logins.

## 2. Docker

An official docker image is available on the Gitlab container registry.
Use it by pulling the following image:

```
registry.gitlab.com/hectorjsmith/fail2ban-prometheus-exporter:latest
```

Use the `:latest` tag to get the most up to date code (less stable) or use one of the version tagged images to use a specific release.
See the [registry page](https://gitlab.com/hectorjsmith/fail2ban-prometheus-exporter/container_registry) for all available tags.

### 2.1. Volumes

The docker image is designed to run by mounting the fail2ban sqlite3 database.
The database should be mounted at: `/app/fail2ban.sqlite3`

The database can be mounted with read-only permissions.

### 2.2. Docker run

Use the following command to run the forwarder as a docker container.

```
docker run -d \
    --name "fail2ban-exporter" \
    -v /var/lib/fail2ban/fail2ban.sqlite3:/app/fail2ban.sqlite3:ro \
    -p "9191:9191"
    registry.gitlab.com/hectorjsmith/fail2ban-prometheus-exporter:latest
```

### 2.3. Docker compose

The following is a simple docker-compose file to run the exporter.

```
version: "2"
services:
  exporter:
    image: registry.gitlab.com/hectorjsmith/fail2ban-prometheus-exporter:latest
    volumes:
    - /var/lib/fail2ban/fail2ban.sqlite3:/app/fail2ban.sqlite3:ro
    ports:
    - "9191:9191"
```

## 3. CLI usage

```
$ fail2ban-prometheus-exporter -h

  -db string
        path to the fail2ban sqlite database
  -port int
        port to use for the metrics server (default 9191)
  -version
        show version info and exit
```

## 4. Metrics

Access exported metrics at `/metrics` (on the provided port).

**Note:** All metric names include the `fail2ban_` prefix to make sure they are unique and easier to find.

Exposed metrics:
* `up` - Returns 1 if the service is up
* `bad_ips` (per jail)
    * A *bad IP* is defined as an IP that has been banned at least once in the past
    * Bad IPs are counted per jail
* `banned_ips` (per jail)
    * A *banned IP* is defined as an IP that is currently banned on the firewall
    * Banned IPs are counted per jail

**Sample**

```
# HELP fail2ban_bad_ips Number of bad IPs stored in the database (per jail).
# TYPE fail2ban_bad_ips gauge
fail2ban_bad_ips{jail="jail1"} 6
fail2ban_bad_ips{jail="jail2"} 8
# HELP fail2ban_banned_ips Number of banned IPs stored in the database (per jail).
# TYPE fail2ban_banned_ips gauge
fail2ban_banned_ips{jail="jail1"} 3
fail2ban_banned_ips{jail="jail2"} 2
# HELP fail2ban_up Was the last fail2ban query successful.
# TYPE fail2ban_up gauge
fail2ban_up 1
```
