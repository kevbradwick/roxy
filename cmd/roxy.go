package main

import (
	"kevbradwick/roxy"
)

func main() {
	cfg := roxy.EnvConfig()
	roxy := roxy.NewRoxy(cfg)
	roxy.Start()
}
