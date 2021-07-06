package rest

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	liveErr "github.com/bal3000/BalStreamerV3/pkg/errors"
	"github.com/bal3000/BalStreamerV3/pkg/livestream"

	"github.com/gorilla/mux"
)

func Handler(l livestream.Service) *mux.Router {
	r := mux.NewRouter()

	// middleware
	r.Use(mux.CORSMethodMiddleware(r))

	// all app routes go here
	s := r.PathPrefix("/api/livestreams").Subrouter()
	s.HandleFunc("/{sportType}/{fromDate}/{toDate}", GetFixtures(l)).Methods(http.MethodGet, http.MethodOptions)
	s.HandleFunc("/{sportType}/{fromDate}/{toDate}/inplay", GetLiveFixtures(l)).Methods(http.MethodGet, http.MethodOptions)
	s.HandleFunc("/{timerId}", GetStreams(l)).Methods(http.MethodGet, http.MethodOptions)

	return r
}

func GetFixtures(l livestream.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		if r.Method == http.MethodOptions {
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		vars := mux.Vars(r)
		sportType := vars["sportType"]
		fromDate := vars["fromDate"]
		toDate := vars["toDate"]

		liveFixtures, err := l.GetLiveFixtures(ctx, sportType, fromDate, toDate, false)
		if err != nil {
			if errors.Is(err, liveErr.StatusErr{StatusCode: 404}) {
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			} else if errors.Is(err, liveErr.StatusErr{StatusCode: 500}) {
				log.Println(err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
			return
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(liveFixtures); err != nil {
			log.Printf("Failed to send json back to client, %v", err)
		}
	}
}

func GetLiveFixtures(l livestream.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		if r.Method == http.MethodOptions {
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		vars := mux.Vars(r)
		sportType := vars["sportType"]
		fromDate := vars["fromDate"]
		toDate := vars["toDate"]

		liveFixtures, err := l.GetLiveFixtures(ctx, sportType, fromDate, toDate, true)
		if err != nil {
			if errors.Is(err, liveErr.StatusErr{StatusCode: 404}) {
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			} else if errors.Is(err, liveErr.StatusErr{StatusCode: 500}) {
				log.Println(err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
			return
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(liveFixtures); err != nil {
			log.Printf("Failed to send json back to client, %v", err)
		}
	}
}

func GetStreams(l livestream.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		if r.Method == http.MethodOptions {
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		vars := mux.Vars(r)
		timerID := vars["timerId"]

		streams, err := l.GetStreams(ctx, timerID)
		if err != nil {
			if errors.Is(err, liveErr.StatusErr{StatusCode: 404}) {
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			} else if errors.Is(err, liveErr.StatusErr{StatusCode: 500}) {
				log.Println(err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
			return
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(streams); err != nil {
			log.Printf("Failed to send json back to client, %v", err)
		}
	}
}
