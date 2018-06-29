# video concat tools
video files concat tools for web

## run / build
Running local server
```sh
sh _run.sh
```

Build docker container
```sh
sh _build.sh
```

## deploy
* https://cloud.google.com/container-optimized-os/docs/how-to/run-container-instance?hl=ja
  * docker-credential-gcr configure-docker
* docker run -it -d -p 80:8080 --name concat-dev asia.gcr.io/planet-pluto-dev/concat-dev

## License
MIT
