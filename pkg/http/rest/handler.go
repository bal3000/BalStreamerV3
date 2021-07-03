package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/bal3000/BalStreamerV3/pkg/livestream"

	"github.com/gorilla/mux"
)

func Handler(l livestream.Service) *mux.Router {
	r := mux.NewRouter()

	// all app routes go here
	s := r.PathPrefix("/api/livestreams").Subrouter()
	s.HandleFunc("/{sportType}/{fromDate}/{toDate}", GetFixtures(l)).Methods(http.MethodGet, http.MethodOptions)
	// s.HandleFunc("/{sportType}/{fromDate}/{toDate}/inplay", GetLiveFixtures).Methods(http.MethodGet, http.MethodOptions)
	// s.HandleFunc("/{timerId}", GetStreams).Methods(http.MethodGet, http.MethodOptions)

	return r
}

func GetFixtures(l livestream.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		if r.Method == http.MethodOptions {
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		vars := mux.Vars(r)
		sportType := vars["sportType"]
		fromDate := vars["fromDate"]
		toDate := vars["toDate"]
		liveStreamURL, apiKey, err := l.GetLiveStreamSettings()
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		url := fmt.Sprintf("%s/%s/%s/%s", liveStreamURL, sportType, fromDate, toDate)
		fixtures := &[]livestream.LiveFixtures{}
		err = callApi(ctx, url, apiKey, fixtures)
		if err != nil {
			log.Println(err)
		}

		if len(*fixtures) == 0 {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(*fixtures); err != nil {
			log.Printf("Failed to send json back to client, %v", err)
		}
	}
}

func callApi(ctx context.Context, url string, apiKey string, body interface{}) error {
	client := &http.Client{}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request to url, %s, err: %w", url, err)
	}
	request.Header.Add("APIKey", apiKey)

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
