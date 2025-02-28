package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Response struct {
	Error string `json:"error,omitempty"`
	Data  any    `json:"data,omitempty"`
}

func sendJSON(w http.ResponseWriter, resp Response, status int) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(resp)
	if err != nil {
		slog.Error("failed to marshal json data", "error", err)
		return
	}

	w.WriteHeader(status)
	if _, err := w.Write(data); err != nil {
		slog.Error("failed to write response to client", "error", err)
		return
	}
}

type application struct {
	Data map[string]string
}

func NewHandler() http.Handler {
	h := chi.NewMux()

	app := application{Data: make(map[string]string)}

	h.Get("/api", HelloWorld(app))

	return h
}

func HelloWorld(app application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sendJSON(w, Response{Data: "Hello World"}, http.StatusOK)
	}
}
