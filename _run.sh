#!/bin/sh

# node process kill
killall node

cd assets
yarn install
yarn run watch &
cd ../

dep ensure -update -v

export ENV=local
export GOOGLE_APPLICATION_CREDENTIALS="./cred/gcs.json"
go run main.go
