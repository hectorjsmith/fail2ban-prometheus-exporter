#/bin/sh

# Print version to logs for debugging purposes
/app/fail2ban-prometheus-exporter --version

socket_path=/var/run/fail2ban/fail2ban.sock
textfile_dir=/app/textfile/
textfile_enabled=false

# Start the exporter (use exec to support graceful shutdown)
# Inspired by: https://akomljen.com/stopping-docker-containers-gracefully/
# Enable textfile metrics if the folder exists (i.e. was mounted by docker)
if [ -d $textfile_dir ]; then
  exec /app/fail2ban-prometheus-exporter \
      --socket "$socket_path" \
      --collector.textfile \
      --collector.textfile.directory "$textfile_dir"
else
  exec /app/fail2ban-prometheus-exporter \
      --socket "$socket_path"
fi
