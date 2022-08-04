FROM golang:alpine as build-env

RUN apk add git

# Build
ENV GOPATH=/go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -v -a -ldflags "-s -w" -o /go/bin/aliyun-exporter .

FROM library/alpine:3.15.0
COPY --from=build-env /go/bin/aliyun-exporter /usr/bin/aliyun-exporter
ENTRYPOINT ["aliyun-exporter"]
