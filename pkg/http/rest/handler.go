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

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		vars := mux.Vars(r)
		sportType := vars["sportType"]
		fromDate := vars["fromDate"]
		toDate := vars["toDate"]

		fixtures := &[]livestream.LiveFixtures{}
		err := l.CallAPI(ctx, fmt.Sprintf("%s/%s/%s", sportType, fromDate, toDate), fixtures)
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
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
