package main

import (
	"context"
	"log"

	"github.com/ninnemana/drudge"
	"github.com/ninnemana/rpc-demo/pkg/service"
)

const (
	tcpAddr = "localhost:8080"
	rpcAddr = "localhost:8081"
)

func main() {
	if err := drudge.Run(context.Background(), drudge.Options{
		Metrics: &drudge.Metrics{
			Prefix:      "tap",
			PullAddress: ":9090",
		},
		BasePath: "/",
		Addr:     tcpAddr,
		RPC: drudge.Endpoint{
			Network: "tcp",
			Addr:    rpcAddr,
		},
		OnRegister: service.Register,
	}); err != nil {
		log.Fatalf("Fell out of serving application: %+v", err)
	}
}
