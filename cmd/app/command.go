package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

func execute(path string) {
	config, err := config.NewConfig(path)
	if err != nil {
		log.Fatal(err)
	}

	storageDB, err := storage.NewStorage(&config.DB)
	if err != nil {
		log.Fatal("error creating storage: ", err)
	}

	driver, err := sqlite3.WithInstance(storageDB.Conn, &sqlite3.Config{})
	if err != nil {
		log.Fatal("error creating driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		"main", driver)
	if err != nil {
		log.Fatalf("migration error: %s", err)
	}

	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Errorf("migration error: %s", err)
	}

	err = storageDB.PrepareAllStatements()
	if err != nil {
		log.Fatalf("error preparing statements: %s", err)
	}

	router, err := router.InitRouter(storageDB)
	if err != nil {
		log.Fatal("router error", err)
	}

	server := &http.Server{
		Addr:              ":" + config.HTTP.Port,
		Handler:           router,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}

	err = subreddit.LoadSubreddits(config.SRFile)
	if err != nil {
		log.Fatal("error loading subreddits: %w", err)
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error("server shutdown:", err)
	}

	if err := storageDB.Close(); err != nil {
		log.Error("error closing database: ", err)
	}

	<-ctx.Done()
	log.Println("timeout of 5 seconds.")
}
