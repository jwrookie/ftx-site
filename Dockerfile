FROM golang:1.17 as builder

ENV GO111MODULE=on
ENV GOPROXY="https://goproxy.cn,direct"
WORKDIR /go/release
COPY . .
RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ftx-site main.go

FROM scratch

WORKDIR /www/ftx-wsite
ENV TIME_ZONE=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TIME_ZONE /etc/localtime && echo $TIME_ZONE > /etc/timezone
COPY --from=builder /go/release/ftx-site /www/ftx-wsite
EXPOSE 8080

ENTRYPOINT ["./ftx-site"]