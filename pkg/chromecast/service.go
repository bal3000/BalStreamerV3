package chromecast

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/bal3000/BalStreamerV3/pkg/errors"
	"github.com/bal3000/BalStreamerV3/pkg/eventbus"
	"github.com/bal3000/BalStreamerV3/pkg/eventbus/events"
	"github.com/bal3000/BalStreamerV3/pkg/storage"
	"github.com/streadway/amqp"
)

var (
	latestEventType = "ChromecastLatestEvent"
	foundEventType  = "ChromecastFoundEvent"
	lostEventType   = "ChromecastLostEvent"
)

type Service interface {
	ListenForChromecasts(routingKey string) error
	GetFoundChromecasts() ([]string, error)
	CastStream(ctx context.Context, routingKey string, c StreamToCast) error
	StopStream(ctx context.Context, routingKey string, c StopPlayingStream) error
}

type service struct {
	eventbus    eventbus.RabbitMQ
	datastore   storage.ChromecastStore
	chromecasts map[string]bool
}

func NewService(e eventbus.RabbitMQ, d storage.ChromecastStore) Service {
	return &service{eventbus: e, datastore: d}
}

func (s *service) ListenForChromecasts(routingKey string) error {
	err := s.eventbus.StartConsumer("chromecast-key", s.processMsgs, 2)
	if err != nil {
		return fmt.Errorf("error consuming rabbit messages: %w", err)
	}

	// send all chromecasts from last refresh to page
	go s.eventbus.SendMessage(routingKey, &events.GetLatestChromecastEvent{MessageType: latestEventType})

	return nil
}

func (s *service) GetFoundChromecasts() ([]string, error) {
	if len(s.chromecasts) == 0 {
		return nil, errors.StatusErr{
			StatusCode: 404,
			Message:    "no chromecasts found",
		}
	}

	casts := make([]string, len(s.chromecasts)-1)
	for k := range s.chromecasts {
		casts = append(casts, k)
	}
	return casts, nil
}

func (s *service) CastStream(ctx context.Context, routingKey string, c StreamToCast) error {
	// Send to chromecast
	cast := &events.StreamToChromecastEvent{
		Chromecast: c.Chromecast,
		StreamURL:  c.StreamURL,
	}

	if err := s.eventbus.SendMessage(routingKey, cast); err != nil {
		return errors.StatusErr{
			StatusCode: 500,
			Message:    err.Error(),
		}
	}

	// save to db
	err := s.datastore.SaveCurrentlyPlaying(ctx, storage.CurrentlyPlaying{
		Fixture:    c.Fixture,
		Chromecast: c.Chromecast,
	})
	if err != nil {
		return errors.StatusErr{
			StatusCode: 500,
			Message:    err.Error(),
		}
	}

	return nil
}

func (s *service) StopStream(ctx context.Context, routingKey string, c StopPlayingStream) error {
	// Send to chromecast
	cast := &events.StopPlayingStreamEvent{
		ChromeCastToStop: c.ChromeCastToStop,
		StopDateTime:     c.StopDateTime,
	}

	if err := s.eventbus.SendMessage(routingKey, cast); err != nil {
		return errors.StatusErr{
			StatusCode: 500,
			Message:    err.Error(),
		}
	}

	// delete from db
	err := s.datastore.DeleteCurrentPlaying(ctx, c.ChromeCastToStop)
	if err != nil {
		return errors.StatusErr{
			StatusCode: 500,
			Message:    err.Error(),
		}
	}

	return nil
}

func (s *service) processMsgs(d amqp.Delivery) bool {
	fmt.Printf("processing message: %s, with type: %s", string(d.Body), d.Type)
	event := new(ChromecastEvent)

	// convert mass transit message
	err := json.Unmarshal(d.Body, event)
	if err != nil {
		log.Println(err)
		return false
	}

	switch event.MessageType {
	case foundEventType:
		s.chromecasts[event.Chromecast] = true
	case lostEventType:
		s.chromecasts[event.Chromecast] = false
	}

	return true
}
