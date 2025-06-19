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
FROM gcr.io/distroless/static-debian12:nonroot
LABEL maintainer="willian.braga@sybogames.com"

USER nonroot

COPY --from=builder /usr/local/bin/moshi-moshi /usr/local/bin/

# The app listens on port 4458 by default via code.
EXPOSE 8080/tcp

CMD ["/usr/local/bin/moshi-moshi"]
