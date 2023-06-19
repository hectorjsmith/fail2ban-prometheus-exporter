FROM debian:buster-slim

# Create main app folder to run from
WORKDIR /app

# Copy compiled binary to release image
# (must build the binary before running docker build)
COPY fail2ban_exporter /app/fail2ban_exporter

ENTRYPOINT ["/app/fail2ban_exporter"]
