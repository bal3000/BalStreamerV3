package chromecast

import (
	"encoding/json"
	"errors"
	"testing"

	customErr "github.com/bal3000/BalStreamerV3/pkg/errors"
	eventMocks "github.com/bal3000/BalStreamerV3/pkg/eventbus/mocks"
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
