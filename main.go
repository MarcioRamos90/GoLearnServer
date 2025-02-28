package main

import (
	"log/slog"
	"net/http"
	"os"
	"server/api"
	"time"
)

func main() {
	if err := run(); err != nil {
		slog.Error("error ")
		os.Exit(1)
	}
	slog.Info("all systems offline")
}

func run() error {
	handler := api.NewHandler()

	s := http.Server{
		ReadTimeout:  time.Second * 5,
		IdleTimeout:  time.Minute * 5,
		WriteTimeout: time.Second * 10,
		Handler:      handler,
		Addr:         ":8080",
	}

	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}
