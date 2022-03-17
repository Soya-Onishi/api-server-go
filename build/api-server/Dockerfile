FROM golang:1.18.0-bullseye as builder

RUN mkdir -p /build/src
WORKDIR /build/src
COPY ../../ ./
RUN go mod download

# クロスコンパイルを有効にする
ARG CGO_ENABLED=0 
ARG GOOS=linux
ARG GOARCH=amd64

# -ldflags '-w -s'によって
# DWARFとシンボルテーブルをバイナリから削除できる
RUN go build -ldflags '-w -s' ./cmd/api-server-go/main.go

FROM alpine
COPY --from=builder /build/src/main /opt/app

ENTRYPOINT ["/opt/app"]