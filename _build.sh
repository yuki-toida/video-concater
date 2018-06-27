#!/bin/sh

# node process kill
killall node

dep ensure

go-assets-builder templates -o assets.go

# docker build & push
docker build -t yukitoida/video-concater:latest .
docker push yukitoida/video-concater

# remove docker image
docker rmi -f yukitoida/video-concater:latest
docker system prune
