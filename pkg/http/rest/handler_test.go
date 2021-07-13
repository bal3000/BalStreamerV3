package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bal3000/BalStreamerV3/pkg/chromecast"
	"github.com/bal3000/BalStreamerV3/pkg/livestream"
	"github.com/bal3000/BalStreamerV3/pkg/livestream/mocks"

	chromeMock "github.com/bal3000/BalStreamerV3/pkg/chromecast/mocks"
)

func TestGetFixtures(t *testing.T) {
	l := mocks.MockService{}

	handler := GetFixtures(l)
	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	url := fmt.Sprintf("%s/Soccer/2021-07-06/2021-07-06", ts.URL)

	res, err := http.Get(url)
	if err != nil {
		t.Errorf("unexpected error occured: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status code %v got %v", http.StatusOK, res.StatusCode)
	}

	fixtures := &[]livestream.LiveFixtures{}

	err = json.NewDecoder(res.Body).Decode(fixtures)
	if err != nil {
		t.Errorf("unexpected error occured: %w", err)
	}

	if len(*fixtures) == 0 {
		t.Errorf("expected fixture count of %v got %v", l.GetFixtureCount(), len(*fixtures))
	}
}

func TestGetStreams(t *testing.T) {
	l := mocks.MockService{}

	handler := GetStreams(l)
	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	url := fmt.Sprintf("%s/1234", ts.URL)

	res, err := http.Get(url)
	if err != nil {
		t.Errorf("unexpected error occured: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status code %v got %v", http.StatusOK, res.StatusCode)
	}

	streams := &livestream.Streams{}

	err = json.NewDecoder(res.Body).Decode(streams)
	if err != nil {
		t.Errorf("unexpected error occured: %w", err)
	}

	if streams.RTMP != l.GetRMTPLink() {
		t.Errorf("expected rtmp link of %v got %v", l.GetRMTPLink(), streams.RTMP)
	}
}

func TestCastStream(t *testing.T) {
	c := chromeMock.MockService{}

	handler := CastStream(c)
	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	url := fmt.Sprintf("%s/api/cast", ts.URL)

	cmd := chromecast.StreamToCast{
		Chromecast: "chromecast 1",
		Fixture:    "test vs testy",
		StreamURL:  "http://test.com",
	}
	json, err := json.Marshal(cmd)
	if err != nil {
		t.Errorf("unexpected error occured: %w", err)
	}

	res, err := http.Post(url, "application/json", bytes.NewBuffer(json))
	if err != nil {
		t.Errorf("unexpected error occured: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent {
		t.Errorf("expected status code %v got %v", http.StatusNoContent, res.StatusCode)
	}
}

func TestStopStream(t *testing.T) {
	c := chromeMock.MockService{}

	handler := StopStream(c)
	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	url := fmt.Sprintf("%s/api/cast", ts.URL)

	cmd := chromecast.StopPlayingStream{
		ChromeCastToStop: "chromecast 1",
		StopDateTime:     time.Now(),
	}
	json, err := json.Marshal(cmd)
	if err != nil {
		t.Errorf("unexpected error occured: %w", err)
	}

	client := &http.Client{}

	req, err := http.NewRequest(http.MethodDelete, url, bytes.NewBuffer(json))
	if err != nil {
		t.Errorf("unexpected error occured: %w", err)
	}
	defer req.Body.Close()

	res, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error occured: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusAccepted {
		t.Errorf("expected status code %v got %v", http.StatusAccepted, res.StatusCode)
	}
}
