package events

import (
	"encoding/json"
	"testing"
	"time"
)

func TestStopTransformMessage(t *testing.T) {
	s := StopPlayingStreamEvent{
		ChromeCastToStop: "test 1",
		StopDateTime:     time.Now(),
	}

	data, tname, err := s.TransformMessage()
	if err != nil {
		t.Errorf("unexpected error occured: %w", err)
	}

	if tname != "StopStreamEvent" {
		t.Errorf("expected type of StopStreamEvent got %s", tname)
	}

	eve := &StopPlayingStreamEvent{}

	err = json.Unmarshal(data, eve)
	if err != nil {
		t.Errorf("unexpected error occured: %w", err)
	}

	if eve.ChromeCastToStop != s.ChromeCastToStop {
		t.Errorf("expected ChromeCastToStop of %s got %s", s.ChromeCastToStop, eve.ChromeCastToStop)
	}
}
