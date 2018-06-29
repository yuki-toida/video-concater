FROM golang:latest as builder
WORKDIR /go/src/github.com/yuki-toida/video-concater/
COPY . .
RUN go get -u github.com/golang/dep/cmd/dep && \
    dep ensure -v && \
    CGO_ENABLED=0 GOOS=linux ENV=dev go build -o app .

FROM alpine:latest
EXPOSE 8080
ENV ENV=dev \
    GOOGLE_APPLICATION_CREDENTIALS="./cred/gcs.json"
RUN apk update && \
    apk upgrade && \
    apk add --no-cache ca-certificates && \
    apk add --no-cache ffmpeg
WORKDIR /opt/app
COPY --from=builder /go/src/github.com/yuki-toida/video-concater/app .
COPY --from=builder /go/src/github.com/yuki-toida/video-concater/index.html .
COPY --from=builder /go/src/github.com/yuki-toida/video-concater/config ./config
COPY --from=builder /go/src/github.com/yuki-toida/video-concater/cred ./cred
CMD ["./app"]
