FROM golang:1.12.7-stretch AS builder
ARG VERSION
RUN apt-get update \
    && apt-get purge -y --auto-remove \
    && apt-get clean -y \
    && rm -rf \
          /var/cache/debconf/* \
          /var/lib/apt/lists/* \
          /var/log/* \
          /tmp/* \
          /var/tmp/*
RUN go get -u golang.org/x/tools/cmd/goimports
ENV SOURCE_DIR /hello-world/
WORKDIR $SOURCE_DIR
COPY . $SOURCE_DIR
RUN make setup \
    && make VERSION=$VERSION


FROM debian:stretch-slim
RUN apt-get update \
    && apt-get install -y --no-install-recommends \
      ca-certificates \
    && apt-get purge -y --auto-remove \
    && apt-get clean -y \
    && rm -rf \
          /var/cache/debconf/* \
          /var/lib/apt/lists/* \
          /var/log/* \
          /tmp/* \
          /var/tmp/*
COPY --from=builder /go/bin/hello-world /app/hello-world

WORKDIR /app
ENV PATH=/app:$PATH

ENTRYPOINT ["/app/hello-world"]
