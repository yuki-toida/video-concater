#!/bin/sh

# node process kill
killall node

cd assets
yarn install
yarn run watch &
cd ../

dep ensure

go-assets-builder templates -o assets.go

ENV=local go run main.go assets.go
