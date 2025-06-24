# syntax=docker/dockerfile:1.11
FROM golang:1.24 AS builder

WORKDIR /tmp/moshi-moshi

COPY ./ /tmp/moshi-moshi

# The CGO_ENABLED will allow the code to run on 'scratch' image.
#RUN go mod download && \
#  CGO_ENABLED=0 go build -o /usr/local/bin/moshi-moshi ./...

# There is no dependency, so let's build it.
RUN CGO_ENABLED=0 go build -o /tmp/moshi-moshi/bin/moshi-moshi ./cmd/moshi-moshi
RUN CGO_ENABLED=0 go build -o /tmp/moshi-moshi/bin/healthcheck ./cmd/healthcheck

## Running image
FROM gcr.io/distroless/static-debian12:nonroot AS runtime

# TODO: Using fixed values, change it to ARGs later on.
LABEL org.opencontainers.image.authors="willian.braga@sybogames.com" \
      org.opencontainers.image.title="Moshi Moshi" \
      org.opencontainers.image.description="A multiservice application for E2E testing." \
      org.opencontainers.image.documentation="https://github.com/sybo-gaming/moshi-moshi/README.md" \
      org.opencontainers.image.licenses="GPL v3" \
      org.opencontainers.image.version="0.0.1" \
      org.opencontainers.image.revision="HEAD" \
      org.opencontainers.image.created="2024-02-12T00:00:00" \
      org.opencontainers.image.source="https://github.com/sybo-gaming/moshi-moshi" \
      org.opencontainers.image.url="https://github.com/sybo-gaming/moshi-moshi"

# TODO: Implement new ideas https://specs.opencontainers.org/image-spec/annotations/
#LABEL vendor="None" \
#      com.example.is-beta= \
#      com.example.is-production="" \
#      com.example.version="0.0.1-beta" \
#      com.example.release-date="2015-02-12"

USER nonroot

COPY --from=builder /tmp/moshi-moshi/bin/* /usr/local/bin/

STOPSIGNAL SIGTERM

HEALTHCHECK --interval=5s --timeout=3s --retries=1 \
  CMD ["/usr/local/bin/moshi-moshi-healthcheck", "http://localhost:8080"]

# The app listens on port 8080 by default via code.
EXPOSE 8080/tcp

# ENTRYPOINT sets the command prefix. The CMD act as a parameter.
ENTRYPOINT ["/usr/local/bin/moshi-moshi"]
