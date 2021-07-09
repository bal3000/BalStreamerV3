package mocks

import (
	"context"

	"github.com/bal3000/BalStreamerV3/pkg/storage"
)

type MockChromecastMongoStore struct{}

func (ds MockChromecastMongoStore) SaveCurrentlyPlaying(ctx context.Context, cp storage.CurrentlyPlaying) error {
	return nil
}

func (ds MockChromecastMongoStore) GetCurrentlyPlaying(ctx context.Context) ([]storage.CurrentlyPlaying, error) {
	return []storage.CurrentlyPlaying{
		{
			Chromecast: "test 1",
			Fixture:    "match 1",
		},
		{
			Chromecast: "test 2",
			Fixture:    "match 2",
		},
	}, nil
}

func (ds MockChromecastMongoStore) DeleteCurrentPlaying(ctx context.Context, chromecast string) error {
	return nil
}
