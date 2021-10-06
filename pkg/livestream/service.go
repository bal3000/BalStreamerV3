package livestream

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/bal3000/BalStreamerV3/pkg/config"
	"github.com/bal3000/BalStreamerV3/pkg/errors"
)

type LiveStreamer interface {
	GetLiveFixtures(ctx context.Context, sportType, fromDate, toDate string, live bool) ([]LiveFixtures, error)
	GetStreams(ctx context.Context, timerID string) (Streams, error)
	FilterLiveFixtures(fixtures []LiveFixtures) ([]LiveFixtures, error)
}

type service struct {
	url    string
	apiKey string
}

func NewService(c config.Configuration) LiveStreamer {
	return service{url: c.LiveStreamURL, apiKey: c.APIKey}
}

func (s service) GetLiveFixtures(ctx context.Context, sportType, fromDate, toDate string, live bool) ([]LiveFixtures, error) {
	fixtures := &[]LiveFixtures{}
	err := s.callAPI(ctx, fmt.Sprintf("%s/%s/%s", sportType, fromDate, toDate), fixtures)
	if err != nil {
		return nil, errors.StatusErr{
			StatusCode: 500,
			Message:    err.Error(),
		}
	}

	if !live {
		if err = checkLen(*fixtures); err != nil {
			return nil, err
		}

		return *fixtures, nil
	}

	lf, err := s.FilterLiveFixtures(*fixtures)
	if err != nil {
		return nil, errors.StatusErr{
			StatusCode: 500,
			Message:    err.Error(),
		}
	}
	if err = checkLen(lf); err != nil {
		return nil, err
	}

	return lf, nil
}

func (s service) GetStreams(ctx context.Context, timerID string) (Streams, error) {
	streams := &Streams{}
	err := s.callAPI(ctx, timerID, streams)
	if err != nil {
		return Streams{}, errors.StatusErr{
			StatusCode: 500,
			Message:    err.Error(),
		}
	}

	return *streams, nil
}

func (s service) FilterLiveFixtures(fixtures []LiveFixtures) ([]LiveFixtures, error) {
	var liveFixtures = []LiveFixtures{}
	for _, fixture := range fixtures {
		if fixture.StateName == "running" {
			liveFixtures = append(liveFixtures, fixture)
			continue
		}

		start, err := parseDate(fixture.UtcStart)
		if err != nil {
			return nil, errors.StatusErr{
				StatusCode: 500,
				Message:    err.Error(),
			}
		}

		end, err := parseDate(fixture.UtcEnd)
		if err != nil {
			return nil, errors.StatusErr{
				StatusCode: 500,
				Message:    err.Error(),
			}
		}

		now := time.Now()

		if (now.Equal(start) || now.After(start)) && now.Before(end) {
			liveFixtures = append(liveFixtures, fixture)
		}
	}

	return liveFixtures, nil
}

func (s service) callAPI(ctx context.Context, path string, body interface{}) error {
	url := fmt.Sprintf("%s/%s", s.url, path)
	client := &http.Client{}
	tctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	request, err := http.NewRequestWithContext(tctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request to url, %s, err: %w", url, err)
	}
	request.Header.Add("APIKey", s.apiKey)

	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("failed to get fixtures from url, %s, err: %w", url, err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusNotFound {
		return fmt.Errorf("url, %s, returned a status code of: %v", url, response.StatusCode)
	}

	if err := json.NewDecoder(response.Body).Decode(body); err != nil {
		return fmt.Errorf("failed to convert JSON, err: %w", err)
	}

	return nil
}

func parseDate(date string) (time.Time, error) {
	layout := "2006-01-02T15:04:05"
	t, err := time.Parse(layout, date)
	if err != nil {
		log.Printf("failed to convert time from live streams, %v", err)
		return time.Time{}, err
	}

	return t, nil
}

func checkLen(f []LiveFixtures) error {
	if len(f) == 0 {
		return errors.StatusErr{
			StatusCode: 404,
			Message:    "no fixtures found",
		}
	}

	return nil
}
