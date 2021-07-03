package livestream

import (
	"errors"

	"github.com/bal3000/BalStreamerV3/pkg/config"
)

type Service interface {
	GetLiveStreamSettings() (string, string, error)
}

type service struct {
	config config.Configuration
}

func NewService(c config.Configuration) Service {
	return &service{config: c}
}

func (s *service) GetLiveStreamSettings() (string, string, error) {
	if s == nil {
		return "", "", errors.New("livestream service is nil")
	}

	return s.config.LiveStreamURL, s.config.APIKey, nil
}
