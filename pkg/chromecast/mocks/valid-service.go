package mocks

import (
	"context"

	"github.com/bal3000/BalStreamerV3/pkg/chromecast"
	"github.com/bal3000/BalStreamerV3/pkg/storage"
)

type MockService struct{}

func (m MockService) ListenForChromecasts(routingKey string) error {
	return nil
}

func (m MockService) GetFoundChromecasts() ([]string, error) {
	return []string{"chromecast 1", "chromecast 2"}, nil
}

func (m MockService) CastStream(ctx context.Context, routingKey string, c chromecast.StreamToCast) error {
	return nil
}

func (m MockService) StopStream(ctx context.Context, routingKey string, c chromecast.StopPlayingStream) error {
	return nil
}

func (m MockService) GetCurrentlyPlayingStream(ctx context.Context) ([]storage.CurrentlyPlaying, error) {
	return []storage.CurrentlyPlaying{
			{
				Chromecast: "chromecast 1",
				Fixture:    "test vs testy",
			},
		},
		nil
}
