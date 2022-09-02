FROM golang as builder
COPY . /go/src/app
WORKDIR /go/src/app
RUN CGO_ENABLED=0 go build -o /build/tsctl ./cmd

FROM registry.access.redhat.com/ubi8/ubi-minimal
COPY --from=builder /build/tsctl /usr/local/bin/tsctl
CMD ["tsctl", "server"]