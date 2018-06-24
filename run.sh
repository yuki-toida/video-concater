#!/bin/sh

# node process kill
killall node

cd assets
yarn install
yarn run dev &
cd ../

dep ensure

go run main.go