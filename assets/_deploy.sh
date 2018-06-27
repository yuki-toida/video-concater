#!/bin/sh

if [ "$1" != "dev" ] && [ "$1" != "prod" ]; then
  echo "error. \$1 undefined env (dev or prod)"
  exit 1
fi

ENV=$1
CACHE_CONTROL='no-cache'
DIST=gs://bucket-tool-$ENV/static/js

echo remove .DS_Store
find . -name '.DS_Store' -type f -ls -delete

echo build $ENV assets
yarn run deploy

echo -n “deploy? [y/n]”
read input

if [ "$input" == "y" ]; then
  echo deploy to $DIST
  gsutil -h "Cache-Control:$CACHE_CONTROL" -m rsync -r -d _build $DIST
  echo remove _build
  rm -rf _build
elif [ "$input" == "n" ]; then
  echo "stop deploy"
else
  echo "input error"
fi
