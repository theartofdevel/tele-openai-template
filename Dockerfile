FROM golang:1.23-alpine AS builder

RUN apk --no-cache add bash git make gcc gettext musl-dev

WORKDIR /usr/local/src/

# cache dependencies
COPY ["app/go.mod", "app/go.sum", "./"]
RUN go mod download

# build application
COPY app/ ./
RUN go build -o ./bin/app cmd/main.go


FROM alpine

RUN apk update && \
    apk upgrade -U && \
    rm -rf /var/cache/*

COPY --from=builder /usr/local/src/bin/app /
COPY configs/config.yml /

CMD ["/app"]
