package main

import (
	"context"
	"log"

	"github.com/ninnemana/drudge"
	"github.com/ninnemana/rpc-demo/pkg/service"
)

const (
	tcpAddr = ":8080"
	rpcAddr = ":8081"
)

var (
	options = drudge.Options{
		BasePath: "/",
		Addr:     tcpAddr,
		RPC: drudge.Endpoint{
			Network: "tcp",
			Addr:    rpcAddr,
		},
		OnRegister: service.Register,
	}
)

func main() {
	if err := drudge.Run(context.Background(), options); err != nil {
		log.Fatalf("Fell out of serving application: %+v", err)
	}
}
