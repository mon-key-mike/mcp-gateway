#syntax=docker/dockerfile:1

ARG GO_VERSION=1.24.5

FROM --platform=${BUILDPLATFORM} golang:${GO_VERSION}-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
COPY vendor/ vendor/
COPY examples/personal-life/server/ examples/personal-life/server/
RUN go build -mod=vendor -o /personal-life-mcp ./examples/personal-life/server

FROM alpine:3.20
RUN adduser -D -u 1000 mcp
USER mcp
COPY --from=builder /personal-life-mcp /usr/local/bin/personal-life-mcp
ENV PERSONAL_MCP_DATA=/data/personal-mcp.json
VOLUME ["/data"]
ENTRYPOINT ["/usr/local/bin/personal-life-mcp"]
