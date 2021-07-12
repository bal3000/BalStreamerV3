package chromecast

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	customErr "github.com/bal3000/BalStreamerV3/pkg/errors"
	eventMocks "github.com/bal3000/BalStreamerV3/pkg/eventbus/mocks"
	"github.com/bal3000/BalStreamerV3/pkg/storage"
	"github.com/bal3000/BalStreamerV3/pkg/storage/mocks"
	"github.com/streadway/amqp"
)

func TestEmptyGetFoundChromecasts(t *testing.T) {
	c := mocks.MockChromecastMongoStore{}
	r := eventMocks.MockRabbitMQ{}

	s := NewService(r, c)

	_, err := s.GetFoundChromecasts()
	if err != nil {
		if !errors.Is(err, customErr.StatusErr{StatusCode: 404}) {
			t.Errorf("expected a 404 to be thrown got %v", err)
		}
	} else {
		t.Errorf("expected an error to be thrown, got nil")
	}
}

func TestGetFoundChromecasts(t *testing.T) {
	msg := ChromecastEvent{
		Chromecast:  "test 1",
		MessageType: "ChromecastFoundEvent",
	}

	b, err := json.Marshal(msg)
	if err != nil {
		t.Errorf("unexpected error occured when marshalling event to json, %s", err.Error())
	}

	c := mocks.MockChromecastMongoStore{}
	r := eventMocks.MockRabbitMQ{
		Msg: amqp.Delivery{
			Body: b,
			Type: "test",
		},
	}

	s := NewService(r, c)

	err = s.ListenForChromecasts("test-router")
	if err != nil {
		t.Errorf("unexpected error occured when listening for chromecasts, %s", err.Error())
	}

	result, err := s.GetFoundChromecasts()
	if err != nil {
		t.Errorf("unexpected error occured, %v", err)
	}

	if len(result) != 1 {
		t.Errorf("expected result length to be 1 got %v", len(result))
	}

	if result[0] != msg.Chromecast {
		t.Errorf("expected result value to be %s got %s", msg.Chromecast, result[0])
	}
}

func TestCastStream(t *testing.T) {
	c := mocks.MockChromecastMongoStore{}
	r := eventMocks.MockRabbitMQ{}

	s := NewService(r, c)

	// add mock call logger

	err := s.CastStream(context.Background(), "test", StreamToCast{
		Chromecast: "test 1",
		Fixture:    "testy",
		StreamURL:  "http://test.com",
	})
	if err != nil {
		t.Errorf("unexpected error occured, %v", err)
	}
}

func TestGetCurrentlyPlayingStream(t *testing.T) {
	m := []storage.CurrentlyPlaying{
		{
			Chromecast: "test 1",
			Fixture:    "match 1",
		},
		{
			Chromecast: "test 2",
			Fixture:    "match 2",
		},
	}

	c := mocks.MockChromecastMongoStore{}
	r := eventMocks.MockRabbitMQ{}

	s := NewService(r, c)

	playing, err := s.GetCurrentlyPlayingStream(context.Background())
	if err != nil {
		t.Errorf("unexpected error occured, %v", err)
	}

	if len(playing) != 2 {
		t.Errorf("expected playing length to be 2 got %v", len(playing))
	}

	for i, p := range playing {
		if p.Chromecast != m[i].Chromecast {
			t.Errorf("expected Chromecast value at index %d to be %s got %s", i, m[i].Chromecast, p.Chromecast)
		}
		if p.Fixture != m[i].Fixture {
			t.Errorf("expected Fixture value at index %d to be %s got %s", i, m[i].Fixture, p.Fixture)
		}
	}
}
