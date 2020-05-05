package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/noptics/services/cmd/registry/data/store"

	"github.com/noptics/golog"
)

func main() {
	if len(os.Getenv("DEBUG")) != 0 {
		golog.Add(golog.StdOut(golog.LEVEL_DEBUG))
	} else {
		golog.Add(golog.StdOut(golog.LEVEL_ERROR))
	}

	defer golog.Finish()

	// Setup the Database Connection
	dbendpoint := os.Getenv("DB_ENDPOINT")

	tablePrefix := os.Getenv("DB_TABLE_PREFIX")
	if len(tablePrefix) == 0 {
		golog.Info("must provide DB_TABLE_PREFIX")
		os.Exit(1)
	}

	dataStore, err := store.New("dynamo", map[string]string{"endpoint": dbendpoint, "prefix": tablePrefix})
	if err != nil {
		golog.Infow("unable to create data store connection", "error", err.Error())
		os.Exit(1)
	}

	// Start the GRPC Server
	gprcPort := os.Getenv("GRPC_PORT")
	if gprcPort == "" {
		gprcPort = "7775"
	}

	errChan := make(chan error)

	gs, err := NewGRPCServer(dataStore, gprcPort, errChan)
	if err != nil {
		golog.Infow("unable to start grpc server", "error", err.Error())
		os.Exit(1)
	}

	// // start the rest server
	// c.RESTPort = os.Getenv("REST_PORT")
	// if c.RESTPort == "" {
	// 	c.RESTPort = "7776"
	// }

	// // We don't support specific address binding right now...
	// c.Host = "0.0.0.0"

	// rs := NewRestServer(dataStore, c.RESTPort, errChan, l, c)

	golog.Info("started")

	// go until told to stop
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

	select {
	case <-sigs:
	case <-errChan:
		golog.Infow("error", "error", err.Error())
	}

	golog.Info("shutting down")
	gs.Stop()

	// err = rs.Stop()
	// if err != nil {
	// 	golog.Infow("error shutting down rest server", "error", err.Error())
	// }

	golog.Info("finished")
}
