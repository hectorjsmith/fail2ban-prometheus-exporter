FROM alpine

# Create main app folder to run from
WORKDIR /app

# Copy compiled binary to release image
# (must build the binary before running docker build)
COPY fail2ban_exporter /app/fail2ban_exporter

# Setup a healthcheck
COPY health /app/health
RUN apk add curl
HEALTHCHECK --interval=10s --timeout=4s --retries=3 CMD /app/health

ENTRYPOINT ["/app/fail2ban_exporter"]
