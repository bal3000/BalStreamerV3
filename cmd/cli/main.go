package main

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"os"

	"github.com/bal3000/BalStreamerV3/pkg/chromecast"
	"github.com/bal3000/BalStreamerV3/pkg/cmd"
	"github.com/bal3000/BalStreamerV3/pkg/config"
	"github.com/bal3000/BalStreamerV3/pkg/eventbus"
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

	cmd := cmd.NewCmdRoot(ls, cs)
	cmd.Execute()
	return nil
}
