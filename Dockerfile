# Using golang:latest instead of alpine because of issues with sqlite3
FROM golang:latest

WORKDIR /app
COPY dist/fail2ban-prometheus-exporter_linux_amd64/fail2ban-prometheus-exporter /app
COPY docker/run.sh /app

RUN chmod +x /app/*

ENTRYPOINT /app/run.sh
