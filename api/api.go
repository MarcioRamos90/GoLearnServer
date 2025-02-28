package api

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"strconv"

	"cloud.google.com/go/firestore"
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
	Data          map[string]string
	ClientService *firestore.Client
	C             context.Context
}

func createClient(ctx context.Context) *firestore.Client {

	client, err := firestore.NewClient(ctx, firestore.DetectProjectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	return client
}

func NewHandler() http.Handler {
	h := chi.NewMux()

	ctx := context.Background()
	clientFirestore := createClient(ctx)

	app := application{Data: make(map[string]string), ClientService: clientFirestore, C: ctx}

	h.Get("/api", GetData(app))
	h.Get("/", HelloApi(app))

	return h
}

func HelloApi(app application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sendJSON(w, Response{Data: "Hello API!"}, http.StatusOK)
	}
}

func GetData(app application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		levelStr := r.URL.Query().Get("level")
		if levelStr == "" {
			levelStr = "0" // default value
		}
		level, err := strconv.ParseInt(levelStr, 10, 64)

		if err != nil {
			slog.Error("Error to parse param", "error", err)
			sendJSON(w, Response{Error: "the param needs to be numeric"}, http.StatusBadRequest)
		}

		q := app.ClientService.Collection("subscriptions").Select().Where("level", "==", level)
		i, err := q.Documents(r.Context()).Next()

		sendJSON(w, Response{Data: i.Exists()}, http.StatusOK)
	}
}
