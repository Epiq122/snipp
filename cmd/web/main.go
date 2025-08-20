package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

type application struct {
	logger *slog.Logger
}

func main() {

	addr := flag.String("addr", ":8080", "http service address")

	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	app := &application{
		logger: logger,
	}

	logger.Info("starting on server", "addr", *addr)

	err := http.ListenAndServe(*addr, app.routes())

	logger.Error(err.Error())
	os.Exit(1)
}
