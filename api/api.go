package api

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/go-chi/chi/v5"
	"google.golang.org/api/iterator"
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
	projectID := firestore.DetectProjectID
	client, err := firestore.NewClient(ctx, projectID)
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

	h.Get("/api", HelloWorld(app))

	return h
}

func HelloWorld(app application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := app.ClientService.Collection("subscriptions").Select().Where("level", "<", 2)

		i, err := q.Documents(r.Context()).Next()

		if err == iterator.Done {
			slog.Error("Failed to get itme", "error", err)
			sendJSON(w, Response{Error: "not found"}, http.StatusNotFound)
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}

		sendJSON(w, Response{Data: i.Data()}, http.StatusOK)
	}
}
