package main

import (
	"bankstats/internal/server"
	"bankstats/internal/service"
	"bankstats/pkg/pkg/cbrclient"
	"flag"
	"log"
	"time"
)

func main() {
	var addr string
	flag.StringVar(&addr, "addr", ":8080", "http listen address")
	flag.Parse()

	cli := cbrclient.New(25 * time.Second)
	svc := &service.Service{
		Client:       cli,
		RequestPause: 200 * time.Millisecond,
	}
	if err := (&server.HTTPServer{Svc: svc}).Run(addr); err != nil {
		log.Fatal(err)
	}
}
