# Roxy

A HTTP reverse proxy with additional features written go Go.

## Run

Run this locally using the Go command

    TARGET=https://www.bbc.co.uk ALLOW_LIST=xxx.xxx.xxx.xxx go run cmd/roxy.go

## Docker image

I've published this to [Docker](https://hub.docker.com/r/kevbradwick/roxy) and you can run it locally with

    docker run --rm -it -e TARGET=https://www.bbc.co.uk -e ALLOW_LIST=x.x.x.x -p 9000:9000 kevbradwick/roxy:latest
