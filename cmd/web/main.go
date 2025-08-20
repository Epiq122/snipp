package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"snippet.robertgleason.ca/internal/models"
)

type application struct {
	logger   *slog.Logger
	snippets *models.SnippetModel
}

func main() {

	addr := flag.String("addr", ":8080", "http service address")

	// db
	dsn := flag.String("dsn", "web:%s@/snippetbox?parseTime=true", "MySQL data source name")

	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Get password from environment
	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		logger.Error("DB_PASSWORD environment variable not set")
		os.Exit(1)
	}
	finalDSN := fmt.Sprintf(*dsn, password)

	db, err := openDB(finalDSN)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	app := &application{
		logger: logger,
		snippets: &models.SnippetModel{
			DB: db,
		},
	}

	logger.Info("starting on server", "addr", *addr)

	err = http.ListenAndServe(*addr, app.routes())

	logger.Error(err.Error())
	os.Exit(1)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}
