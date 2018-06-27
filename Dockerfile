FROM golang:alpine as builder
WORKDIR /go/src/github.com/yuki-toida/video-concater/
COPY . .
RUN GOOS=linux ENV=dev go build -o app .

FROM alpine:latest
EXPOSE 8080
ENV ENV=dev
RUN apk update && \
    apk upgrade && \
    apk --no-cache add ffmpeg
WORKDIR /opt/app
COPY --from=builder /go/src/github.com/yuki-toida/video-concater/app .
COPY --from=builder /go/src/github.com/yuki-toida/video-concater/config ./config
CMD ["./app"]
