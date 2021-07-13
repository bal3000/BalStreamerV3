package main

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/bal3000/BalStreamerV3/pkg/chromecast"
	"github.com/bal3000/BalStreamerV3/pkg/config"
	"github.com/bal3000/BalStreamerV3/pkg/eventbus"
	"github.com/bal3000/BalStreamerV3/pkg/http/rest"
	"github.com/bal3000/BalStreamerV3/pkg/livestream"
	"github.com/bal3000/BalStreamerV3/pkg/storage/mongo"
)

//go:embed config.json
var jf []byte

var configuration config.Configuration

func init() {
	configuration = config.ReadConfig(bytes.NewBuffer(jf))
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "startup error: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	//setup rabbit
	rabbit, closer, err := eventbus.NewRabbitMQ(configuration)
	if err != nil {
		return err
	}
	defer closer()

	// setup chromecast db
	mongo, dbCloser, err := mongo.NewChromecastMongoStore(context.Background(), configuration.ConnectionString)
	if err != nil {
		return err
	}
	defer dbCloser()

	// services
	ls := livestream.NewService(configuration)
	cs := chromecast.NewService(rabbit, mongo)

	// routes and middleware
	router := rest.Handler(ls, cs)

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
