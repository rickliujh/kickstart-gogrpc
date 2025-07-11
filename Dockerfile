# Type of binary that will be built, one of [prod, dev]
ARG BUILD_TYPE=dev
# Specify name of the binary being built
ARG APP_NAME=kickstart-go
# Specify the image that will run the binary
ARG RUN_IMG=scratch



FROM golang:alpine AS builder

ARG BUILD_TYPE
ARG APP_NAME

WORKDIR /go/src/app
COPY . .

RUN apk add git make

RUN make build-${BUILD_TYPE}



FROM ${RUN_IMG}

ARG APP_NAME

COPY --from=builder /go/src/app/bin/${APP_NAME} /app

# the tls certificates:
# NB: this pulls directly from the upstream image, which already has ca-certificates:
COPY --from=alpine:latest /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/app"]
