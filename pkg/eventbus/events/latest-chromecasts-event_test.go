package events

import (
	"encoding/json"
	"testing"
)

func TestLatestTransformMessage(t *testing.T) {
	s := GetLatestChromecastEvent{
		MessageType: "test 1",
	}

	data, tname, err := s.TransformMessage()
	if err != nil {
		t.Errorf("unexpected error occured: %w", err)
	}

	if tname != "ChromecastLatestEvent" {
		t.Errorf("expected type of StopStreamEvent got %s", tname)
	}

	eve := &GetLatestChromecastEvent{}

	err = json.Unmarshal(data, eve)
	if err != nil {
		t.Errorf("unexpected error occured: %w", err)
	}

	if eve.MessageType != s.MessageType {
		t.Errorf("expected MessageType of %s got %s", s.MessageType, eve.MessageType)
	}
}
