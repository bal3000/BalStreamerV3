package storage

import "context"

type ChromecastStore interface {
	SaveCurrentlyPlaying(ctx context.Context, cp CurrentlyPlaying) error
	GetCurrentlyPlaying(ctx context.Context) ([]CurrentlyPlaying, error)
	DeleteCurrentPlaying(ctx context.Context, chromecast string) error
}
