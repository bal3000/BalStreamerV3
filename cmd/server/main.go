package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/bal3000/BalStreamerV3/pkg/config"
	"github.com/bal3000/BalStreamerV3/pkg/http/rest"
	"github.com/bal3000/BalStreamerV3/pkg/livestream"
)

var configuration config.Configuration

func init() {
	file, _ := os.Open("./config.json")
	defer file.Close()
	configuration = config.ReadConfig(file)
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "startup error: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	// services
	ls := livestream.NewService(configuration)

	// routes and middleware
	router := rest.Handler(ls)

	// start the server
	srv := &http.Server{
		Addr:         "0.0.0.0:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
		Handler:      router,
	}

	go func() {
		log.Println("started server on port 8080")
		if err := srv.ListenAndServe(); err != nil {
			panic(err)
		}
	}()

	// listen for a ctrl+c and graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		return err
	}

	log.Println("shutting down")
	os.Exit(0)
	return nil
}
