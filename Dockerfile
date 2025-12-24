FROM golang:1.25-alpine as builder
ARG MODULE

WORKDIR /app
COPY . /app
RUN apk add gcc musl-dev
RUN go install -tags musl

ENV GOOS=linux GOARCH=amd64 GOPROXY=https://proxy.golang.org GOEXPERIMENT=jsonv2
RUN go build -tags musl -ldflags="-w -s" -o new_project main.go
RUN chmod +x new_project

FROM alpine
RUN apk add curl
WORKDIR /app
COPY --from=builder /app/new_project .
COPY ./config /app/config
ENTRYPOINT ["/app/new_project"]
