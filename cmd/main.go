package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/ninnemana/drudge"
	"github.com/ninnemana/rpc-demo/pkg/vinyltap"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

const (
	tcpAddr = ":8080"
	rpcAddr = ":8081"
)

var (
	options = drudge.Options{
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
	}
)

func main() {
	if err := drudge.Run(context.Background(), options); err != nil {
		log.Fatalf("Fell out of serving application: %+v", err)
	}
}

type Service struct {
	albums map[int32]*vinyltap.Album
	sync.RWMutex
}

func Register(server *grpc.Server) error {
	vinyltap.RegisterTapServer(server, &Service{
		albums: map[int32]*vinyltap.Album{},
	})
	return nil
}

func (s *Service) GetAlbum(a *vinyltap.Album, srv vinyltap.Tap_GetAlbumServer) error {
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

func (s *Service) Set(ctx context.Context, a *vinyltap.Album) (*vinyltap.Album, error) {
	if _, ok := s.albums[a.GetId()]; ok {
		return nil, errors.Errorf("provided album '%d' exists", a.GetId())
	}
	s.Lock()
	s.albums[a.GetId()] = a
	s.Unlock()

	return s.albums[a.GetId()], nil
}

func sinceInMilliseconds(startTime time.Time) float64 {
	return float64(time.Since(startTime).Nanoseconds()) / 1e6
}
