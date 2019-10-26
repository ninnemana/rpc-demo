package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"sync"

	"cloud.google.com/go/firestore"
	"github.com/ninnemana/drudge"
	"github.com/pkg/errors"
	"github.com/uber/jaeger-client-go/config"
	option "google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ninnemana/rpc-demo/pkg/vinyltap"
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
		SwaggerDir:    "openapi",
		Handlers:      []drudge.Handler{vinyltap.RegisterTapHandler},
		OnRegister:    Register,
		TraceExporter: drudge.Jaeger,
	}
)

func main() {
	cfg, err := config.FromEnv()
	if err != nil {
		log.Fatalf("Failed to read tracing config: %v", err)
	}

	cfg.ServiceName = "rpc-demo"
	options.TraceConfig = cfg

	if err := drudge.Run(context.Background(), options); err != nil {
		log.Fatalf("Fell out of serving application: %+v", err)
	}
}

type Service struct {
	albums map[int32]*vinyltap.Album
	client *firestore.Client
	sync.RWMutex
}

func Register(server *grpc.Server) error {
	options := []option.ClientOption{}

	if os.Getenv("GCP_SVC_ACCOUNT") != "" {
		js, err := base64.StdEncoding.DecodeString(os.Getenv("GCP_SVC_ACCOUNT"))
		if err != nil {
			return errors.WithMessage(err, "failed to decode service account")
		}

		options = append(options, option.WithCredentialsJSON(js))
	}

	client, err := firestore.NewClient(context.Background(), os.Getenv("GCP_PROJECT_ID"), options...)
	if err != nil {
		return errors.WithMessage(err, "failed to create Firestore client")
	}

	service := &Service{
		albums: map[int32]*vinyltap.Album{},
		client: client,
	}

	vinyltap.RegisterTapServer(server, service)

	return nil
}

func (s *Service) GetAlbum(a *vinyltap.Album, srv vinyltap.Tap_GetAlbumServer) error {
	snap, err := s.client.
		Collection("album").
		Doc(fmt.Sprintf("%d", a.GetId())).
		Get(srv.Context())

	switch {
	case err == nil:
	case status.Code(err) == codes.NotFound:
		return err
	default:
		return status.Error(
			codes.Internal,
			errors.WithMessage(err, "failed to find album").Error(),
		)
	}

	album := &vinyltap.Album{}
	if err := snap.DataTo(album); err != nil {
		return status.Error(
			codes.Internal,
			errors.WithMessage(err, "failed to read album").Error(),
		)
	}

	if err := srv.Send(album); err != nil {
		return status.Error(
			codes.Unavailable,
			errors.WithMessage(err, "failed to send album").Error(),
		)
	}

	return nil
}

func (s *Service) Set(ctx context.Context, a *vinyltap.Album) (*vinyltap.Album, error) {
	if a.GetId() > 0 {
		_, err := s.client.Collection("album").Doc(fmt.Sprintf("%d", a.GetId())).Get(ctx)

		switch {
		case err == nil:
			return nil, status.Errorf(
				codes.AlreadyExists,
				errors.WithMessage(err, "the record we are setting already exists").Error(),
			)
		case status.Code(err) == codes.NotFound:
		default:
			return nil, status.Errorf(
				codes.Internal,
				errors.WithMessage(err, "failed to check for existing record").Error(),
			)
		}
	}

	_, err := s.client.Collection("album").Doc(
		fmt.Sprintf("%d", a.GetId()),
	).Set(ctx, a)
	if err != nil {
		return nil, status.Error(
			codes.Internal,
			errors.WithMessage(err, "failed to write album").Error(),
		)
	}

	return a, nil
}
