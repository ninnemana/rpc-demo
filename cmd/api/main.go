package main

import (
	"context"
	"log"
	"sync"

	"github.com/ninnemana/drudge"
	"github.com/ninnemana/drudge/telemetry"
	"github.com/ninnemana/rpc-demo/pkg/vinyltappb"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
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
		OnRegister:    Register,
		TraceExporter: telemetry.Jaeger,
		TraceConfig: telemetry.JaegerConfig{
			ServiceName: "rpc-demo",
		},
	}
)

func main() {
	if err := drudge.Run(context.Background(), options); err != nil {
		log.Fatalf("Fell out of serving application: %+v", err)
	}
}

type Service struct {
	albums map[int32]*vinyltappb.Album
	sync.RWMutex
}

func Register(server *grpc.Server) error {
	vinyltap.RegisterTapServer(server, &Service{
		albums: map[int32]*vinyltap.Album{},
	})
	return nil
}

func (s *Service) GetAlbum(a *vinyltappb.Album, srv vinyltappb.Tap_GetAlbumServer) error {
	for k := range s.albums {
		if k != a.GetId() {
			continue
		}

		if err := srv.Send(s.albums[k]); err != nil {
			return errors.WithMessage(err, "failed to send album over TCP connection")
		}
	}

	return nil
}

func (s *Service) Set(ctx context.Context, a *vinyltappb.Album) (*vinyltappb.Album, error) {
	if _, ok := s.albums[a.GetId()]; ok {
		return nil, errors.Errorf("provided album '%d' exists", a.GetId())
	}
	s.Lock()
	s.albums[a.GetId()] = a
	s.Unlock()

	return s.albums[a.GetId()], nil
}
