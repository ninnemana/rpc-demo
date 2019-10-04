package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ninnemana/drudge/telemetry"
	"github.com/ninnemana/rpc-demo/pkg/vinyltappb"
	"go.opencensus.io/stats"
	"go.opencensus.io/tag"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type Service struct {
	albums map[int32]*vinyltappb.Album
	sync.RWMutex
}

func Register(server *grpc.Server) error {
	vinyltappb.RegisterTapServer(server, &Service{})
	return nil
}

var (
	GetLatencyMs = telemetry.Float64Measure(
		"getAlbum/latency",
		"The latency in milliseconds per GetAlbum request",
		"ms",
		[]tag.Key{
			KeyID,
		},
		telemetry.LatencyDistribution,
	)
	SetLatencyMs = telemetry.Float64Measure(
		"set/latency",
		"The latency in milliseconds per Set request",
		"ms",
		[]tag.Key{
			KeyID,
		},
		telemetry.LatencyDistribution,
	)

	KeyID, _ = tag.NewKey("id")
)

func (s *Service) GetAlbum(a *vinyltappb.Album, srv vinyltappb.Tap_GetAlbumServer) error {
	span, ctx := opentracing.StartSpanFromContext(srv.Context(), "vinyltappb.GetAlbum")
	span.SetTag("id", a.GetId())
	ctx, _ = tag.New(ctx, tag.Insert(
		KeyID,
		fmt.Sprintf("%d", a.GetId()),
	))
	defer func(start time.Time) {
		span.Finish()
		stats.Record(ctx, GetLatencyMs.M(sinceInMilliseconds(start)))
	}(time.Now())

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
	span, ctx := opentracing.StartSpanFromContext(ctx, "vinyltappb.Set")
	span.SetTag("id", a.GetId())
	ctx, _ = tag.New(ctx, tag.Insert(
		KeyID,
		fmt.Sprintf("%d", a.GetId()),
	))
	defer func(start time.Time) {
		span.Finish()
		stats.Record(ctx, SetLatencyMs.M(sinceInMilliseconds(start)))
	}(time.Now())

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
