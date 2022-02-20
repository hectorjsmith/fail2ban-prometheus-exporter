# Using golang:latest instead of alpine because of issues with sqlite3
FROM golang:latest AS build

# Create build folder to compile tool
WORKDIR /build

# Copy source files to build folder and link to the /go folder
COPY . /build
RUN ln -s /go/src/ /build/src

# Compile the tool using a Make command
RUN make build/docker


FROM debian:buster-slim

# Create main app folder to run from
WORKDIR /app

# Copy compiled binary to release image
COPY --from=build /build/src/fail2ban_exporter /app/fail2ban_exporter

ENTRYPOINT ["/app/fail2ban_exporter"]
