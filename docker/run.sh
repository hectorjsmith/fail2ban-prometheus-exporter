#/bin/sh

# Print version to logs for debugging purposes
/app/fail2ban-prometheus-exporter -version

# Start the exporter (use exec to support graceful shutdown)
# Inspired by: https://akomljen.com/stopping-docker-containers-gracefully/
exec /app/fail2ban-prometheus-exporter \
    -db /app/fail2ban.sqlite3
