FROM golang as builder
COPY . /go/src/app
WORKDIR /go/src/app
ARG CLI_VERSION
RUN CGO_ENABLED=0 go build -o /build/tsctl -ldflags "-X main.cliVersion=$CLI_VERSION" ./cmd

FROM registry.access.redhat.com/ubi8/ubi-minimal
COPY --from=builder /build/tsctl /usr/local/bin/tsctl
CMD ["tsctl", "server"]