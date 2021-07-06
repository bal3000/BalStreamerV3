package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bal3000/BalStreamerV3/pkg/livestream"
	"github.com/bal3000/BalStreamerV3/pkg/livestream/mocks"
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
