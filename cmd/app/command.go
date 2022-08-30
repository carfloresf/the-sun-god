package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	migrate "github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"

	"github.com/carfloresf/the-sun-god/config"
	"github.com/carfloresf/the-sun-god/internal/router"
	"github.com/carfloresf/the-sun-god/pkg/storage"
	"github.com/carfloresf/the-sun-god/pkg/subreddit"
)

func execute(path string) error {
	config, err := config.NewConfig(path)
	if err != nil {
		log.Errorln("error loading config: ", err)
		return err
	}

	storageDB, err := storage.NewStorage(&config.DB)
	if err != nil {
		log.Errorf("error creating storage: %s", err)
		return err
	}

	driver, err := sqlite3.WithInstance(storageDB.Conn, &sqlite3.Config{})
	if err != nil {
		log.Errorf("error creating driver: %s", err)
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		"main", driver)
	if err != nil {
		log.Errorf("migration error: %s", err)
		return err
	}

	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Errorf("migration error: %s", err)
		return err
	}

	err = storageDB.PrepareAllStatements()
	if err != nil {
		log.Errorf("error preparing statements: %s", err)
		return err
	}

	router, err := router.InitRouter(storageDB)
	if err != nil {
		log.Errorf("router error: %s", err)
	}

	server := &http.Server{
		Addr:              ":" + config.HTTP.Port,
		Handler:           router,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}

	err = subreddit.LoadSubreddits(config.SRFile)
	if err != nil {
		log.Errorf("error loading subreddits: %s", err)
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Errorf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Errorf("server shutdown: %s", err)
	}

	if err := storageDB.Close(); err != nil {
		log.Errorf("error closing database: %s", err)
	}

	return nil
}
