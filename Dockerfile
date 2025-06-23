# syntax=docker/dockerfile:1.11
FROM golang:1.24 AS builder

COPY ./cmd /tmp/moshi-moshi
WORKDIR /tmp/moshi-moshi
COPY go.* .

# The CGO_ENABLED will allow the code to run on 'scratch' image.
#RUN go mod download && \
#  CGO_ENABLED=0 go build -o /usr/local/bin/moshi-moshi ./...

# There is no dependency, so let's build it.
RUN CGO_ENABLED=0 go build -o /usr/local/bin/moshi-moshi ./...

## Running image
FROM gcr.io/distroless/static-debian12:nonroot AS runtime
LABEL org.opencontainers.image.authors="willian.braga@sybogames.com"

# TODO: Implement new ideas https://specs.opencontainers.org/image-spec/annotations/
#LABEL vendor="None" \
#      com.example.is-beta= \
#      com.example.is-production="" \
#      com.example.version="0.0.1-beta" \
#      com.example.release-date="2015-02-12"

USER nonroot

ONBUILD COPY --from=builder /usr/local/bin/moshi-moshi /usr/local/bin/

STOPSIGNAL SIGTERM

#HEALTHCHECK --interval=5s --timeout=3s --retries=1 \
#  CMD /usr/local/bin/moshi-moshi-healthcheck http://localhost:8080

# The app listens on port 8080 by default via code.
EXPOSE 8080/tcp

# ENTRYPOINT sets the command prefix. The CMD act as a parameter.
ENTRYPOINT ["/usr/local/bin/moshi-moshi"]
