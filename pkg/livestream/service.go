package livestream

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/bal3000/BalStreamerV3/pkg/config"
)

type Service interface {
	CallAPI(ctx context.Context, path string, body interface{}) error
}

type service struct {
	url    string
	apiKey string
}

func NewService(c config.Configuration) Service {
	return service{url: c.LiveStreamURL, apiKey: c.APIKey}
}

func (s service) CallAPI(ctx context.Context, path string, body interface{}) error {
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
