#/bin/sh

# Print version to logs for debugging purposes
/app/fail2ban-prometheus-exporter -version

db_path=/app/fail2ban.sqlite3
socket_path=/var/run/fail2ban/fail2ban.sock
textfile_dir=/app/textfile/
textfile_enabled=false

# Blank out the file paths if they do not exist - a hacky way to only use these files if they were mounted into the container.
if [ ! -f "$db_path" ]; then
    db_path=""
fi
if [ ! -S "$socket_path" ]; then
    socket_path=""
fi
if [ -d $textfile_dir ]; then
    textfile_enabled=true
fi

# Start the exporter (use exec to support graceful shutdown)
# Inspired by: https://akomljen.com/stopping-docker-containers-gracefully/
exec /app/fail2ban-prometheus-exporter \
    -db "$db_path" \
    -socket "$socket_path" \
    -collector.textfile=$textfile_enabled \
    -collector.textfile.directory="$textfile_dir"
