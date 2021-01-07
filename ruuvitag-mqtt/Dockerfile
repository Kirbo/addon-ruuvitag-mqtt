ARG BUILD_FROM=alpine
FROM golang:1.15.2-alpine AS builder

# Local build:
# docker build -t ruuvitag-mqtt . && docker run --rm -v $PWD/data/options.example.json:/data/options.json ruuvitag-mqtt

ARG BUILD_ARCH
ARG CLI_VERSION

# Set necessary environment variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0

# Move to working directory /build
WORKDIR /build

# Copy static assets into the container
COPY go.mod go.sum ./

# Download dependencies using go mod
RUN go mod download

# Copy the code into the container
COPY cmd ./cmd
COPY internal ./internal
COPY entrypoint.sh ./

# Build the application
RUN go build -o main cmd/main.go

# Reduce end image size by dropping build phases out
# hadolint ignore=DL3006
# FROM alpine
FROM $BUILD_FROM

ARG BUILD_DATE
ARG BUILD_REF
ARG BUILD_VERSION=0.0.1

RUN apk add bluez

# Move to /ruuvitag directory as the place for resulting binary folder
WORKDIR /ruuvitag

# Copy binary from build to main folder
COPY --from=builder /build/entrypoint.sh ./
COPY --from=builder /build/main ./

# Command to run when starting the container
CMD ["/ruuvitag/entrypoint.sh"]

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
