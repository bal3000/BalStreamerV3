package middleware

import (
	"io"
	"io/ioutil"
	"net/http"
)

func DrainAndClose(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
			_, _ = io.Copy(ioutil.Discard, r.Body)
			_ = r.Body.Close()
		},
	)
}
