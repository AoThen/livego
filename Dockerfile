FROM golang:alpine as builder
WORKDIR /app
ENV GOPROXY https://goproxy.io
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -a -installsuffix cgo -o livego . 

FROM alpine:latest
RUN mkdir -p /app/config && apk add --no-cache tzdata \
&& cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime  \
&& echo Asia/Shanghai > /etc/timezone \
&& apk del tzdata
WORKDIR /app
ENV RTMP_PORT 1935
ENV HTTP_FLV_PORT 7001
ENV HLS_PORT 7002
ENV HTTP_OPERATION_PORT 8090
COPY --from=builder /app/livego .
EXPOSE ${RTMP_PORT}
EXPOSE ${HTTP_FLV_PORT}
EXPOSE ${HLS_PORT}
EXPOSE ${HTTP_OPERATION_PORT}
ENTRYPOINT ["./livego"]
