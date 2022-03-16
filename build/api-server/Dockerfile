FROM golang:1.18.0-bullseye as builder

WORKDIR /build
COPY ../../go.mod ../../go.sum ./
RUN go mod download
COPY ../../ ./

# クロスコンパイルを有効にする
# ARG CGO_ENABLED=0 
# ARG GOOS=linux
# ARG GOARCH=amd64

# -ldflags '-w -s'によって
# DWARFとシンボルテーブルをバイナリから削除できる
RUN go build -ldflags '-w -s' ./cmd/api-server-go/main.go

FROM alpine
COPY --from=builder /build/main /opt/app

ENTRYPOINT ["opt/app/main"]