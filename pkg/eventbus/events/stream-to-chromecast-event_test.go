package events

import (
	"encoding/json"
	"testing"
)

func TestStreamTransformMessage(t *testing.T) {
	s := StreamToChromecastEvent{
		Chromecast: "test 1",
		StreamURL:  "test",
	}

	data, tname, err := s.TransformMessage()
	if err != nil {
		t.Errorf("unexpected error occured: %w", err)
	}

	if tname != "PlayStreamEvent" {
		t.Errorf("expected type of PlayStreamEvent got %s", tname)
	}

	eve := &StreamToChromecastEvent{}

	err = json.Unmarshal(data, eve)
	if err != nil {
		t.Errorf("unexpected error occured: %w", err)
	}

	if eve.Chromecast != s.Chromecast {
		t.Errorf("expected Chromecast of %s got %s", s.Chromecast, eve.Chromecast)
	}
}
