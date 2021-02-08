# Using golang:latest instead of alpine because of issues with sqlite3
FROM golang:latest

# Create build folder to compile tool
WORKDIR /build

# Copy source files to build folder and link to the /go folder
COPY . /build
RUN ln -s /go/src/ /build/src

# Compile the tool using a Make command
RUN make build/docker

# Create main app folder to run from
WORKDIR /app

# Move compiled binary to app folder and delete build folder
RUN mv /build/src/exporter /app/fail2ban-prometheus-exporter
RUN rm -rf /build

# Copy init script into main app folder and set as entry point
COPY docker/run.sh /app/
RUN chmod +x /app/*
ENTRYPOINT /app/run.sh
