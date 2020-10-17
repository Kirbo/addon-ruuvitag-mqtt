FROM golang:1.15.2-alpine AS builder

# Set necessary environment variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Move to working directory /build
WORKDIR /build

# Copy static assets into the container
COPY go.mod go.sum ./

# Download dependencies using go mod
RUN go mod download

# Copy the code into the container
COPY cmd ./cmd
COPY internal ./internal

# Build the application
RUN go build -o main cmd/main.go

# Reduce end image size by dropping build phases out
FROM alpine:3.12.0

# Move to /dist directory as the place for resulting binary folder
WORKDIR /dist

# Command to run when starting the container
CMD ["/dist/main"]

# Build arugments
ARG BUILD_DATE
ARG BUILD_REF
ARG BUILD_VERSION

# Labels
LABEL \
    io.hass.name="RuuviTag into MQTT" \
    io.hass.description="Broadcast RuuviTag sensors into MQTT" \
    io.hass.arch="${BUILD_ARCH}" \
    io.hass.type="addon" \
    io.hass.version=${BUILD_VERSION} \
    maintainer="Kimmo Saari <kirbo@kirbo-designs.com>" \
    org.label-schema.description="Broadcast RuuviTag sensors into MQTT" \
    org.label-schema.build-date=${BUILD_DATE} \
    org.label-schema.name="RuuviTag into MQTT" \
    org.label-schema.schema-version="1.0" \
    org.label-schema.usage="https://gitlab.com/kirbo/addon-ruuvitag-mqtt/-/blob/master/README.md" \
    org.label-schema.vcs-ref=${BUILD_REF} \
    org.label-schema.vcs-url="https://gitlab.com/kirbo/addon-ruuvitag-mqtt/" \
    org.label-schema.vendor="HomeAssistant add-ons by Kimmo Saari"
