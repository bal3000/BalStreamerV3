package rest

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/bal3000/BalStreamerV3/pkg/chromecast"
	liveErr "github.com/bal3000/BalStreamerV3/pkg/errors"
	"github.com/bal3000/BalStreamerV3/pkg/http/middleware"
	"github.com/bal3000/BalStreamerV3/pkg/livestream"

	"github.com/gorilla/mux"
)

const routingKey string = "chromecast-key"

func Handler(l livestream.Service, c chromecast.Service) *mux.Router {
	r := mux.NewRouter()

	// middleware
	r.Use(mux.CORSMethodMiddleware(r))
	r.Use(middleware.DrainAndClose)

	// all app routes go here
	s := r.PathPrefix("/api/livestreams").Subrouter()
	s.HandleFunc("/{sportType}/{fromDate}/{toDate}", GetFixtures(l)).Methods(http.MethodGet, http.MethodOptions)
	s.HandleFunc("/{sportType}/{fromDate}/{toDate}/inplay", GetLiveFixtures(l)).Methods(http.MethodGet, http.MethodOptions)
	s.HandleFunc("/{timerId}", GetStreams(l)).Methods(http.MethodGet, http.MethodOptions)

	sc := r.PathPrefix("/api/cast").Subrouter()
	sc.HandleFunc("", CastStream(c)).Methods(http.MethodPost, http.MethodOptions)
	sc.HandleFunc("", StopStream(c)).Methods(http.MethodDelete, http.MethodOptions)

	r.HandleFunc("/api/chromecasts", GetChromecasts(c)).Methods(http.MethodGet)
	r.HandleFunc("/api/currentplaying", GetCurrentlyPlayingStream(c)).Methods(http.MethodGet)

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

func GetChromecasts(c chromecast.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		if r.Method == http.MethodOptions {
			return
		}

		chromecasts, err := c.GetFoundChromecasts()
		if err != nil {
			if errors.Is(err, liveErr.StatusErr{StatusCode: 404}) {
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			} else {
				log.Println(err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
			return
		}

		if len(chromecasts) == 0 {
			http.NotFound(w, r)
			return
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(chromecasts); err != nil {
			log.Printf("Failed to send json back to client, %v", err)
		}
	}
}

func GetCurrentlyPlayingStream(c chromecast.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		if r.Method == http.MethodOptions {
			return
		}

		playing, err := c.GetCurrentlyPlayingStream(r.Context())
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
		if err := json.NewEncoder(w).Encode(playing); err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}
}

// CastStream - streams given data to given chromecast
func CastStream(c chromecast.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "DELETE,HEAD,OPTIONS,POST,PUT")
		w.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers")

		if r.Method == http.MethodOptions {
			return
		}

		castCommand := new(chromecast.StreamToCast)

		if err := json.NewDecoder(r.Body).Decode(castCommand); err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}

		// Send to chromecast
		err := c.CastStream(r.Context(), routingKey, *castCommand)
		if err != nil {
			if errors.Is(err, liveErr.StatusErr{StatusCode: 404}) {
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			} else if errors.Is(err, liveErr.StatusErr{StatusCode: 500}) {
				log.Println(err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// StopStream endpoint sends the command to stop the stream on the given chromecast
func StopStream(c chromecast.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "DELETE,HEAD,OPTIONS,POST,PUT")
		w.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers")

		log.Println("method", r.Method)
		if r.Method == http.MethodOptions {
			return
		}

		stopStreamCommand := new(chromecast.StopPlayingStream)

		if err := json.NewDecoder(r.Body).Decode(stopStreamCommand); err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}

		// Send to chromecast
		err := c.StopStream(r.Context(), routingKey, *stopStreamCommand)
		if err != nil {
			if errors.Is(err, liveErr.StatusErr{StatusCode: 404}) {
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			} else if errors.Is(err, liveErr.StatusErr{StatusCode: 500}) {
				log.Println(err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
			return
		}

		w.WriteHeader(http.StatusAccepted)
	}
}
