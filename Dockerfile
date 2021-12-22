FROM golang:1.15.3-alpine3.12 as build
RUN apk add make git
ADD . /go/src/icode.baidu.com/sim-exporter
WORKDIR /go/src/icode.baidu.com/sim-exporter
ENV GOPROXY=https://goproxy.cn
RUN make

FROM alpine:3.12
COPY --from=build /go/src/icode.baidu.com/sim-exporter/bin/sim-exporter /usr/bin/sim-exporter
RUN chmod +x /usr/bin/sim-exporter
ENTRYPOINT ["/usr/bin/sim-exporter"]
