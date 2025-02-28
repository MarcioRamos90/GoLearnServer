package api

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewHandler() http.Handler {
	h := chi.NewMux()

	ctx := context.Background()
	clientFirestore := createClient(ctx)

	app := application{Data: make(map[string]string), FirestoreClient: clientFirestore}

	h.Get("/", HelloApi(app))
	h.Get("/api", GetData(app))

	return h
}

func HelloApi(app application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sendJSON(w, Response{Data: "Hello API!"}, http.StatusOK)
	}
}

func GetData(app application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		level := extracIntegerFromQueryParam(w, r, "level")
		q := app.FirestoreClient.Collection("subscriptions").Select().Where("level", "==", level)
		i, err := q.Documents(r.Context()).GetAll()

		if err != nil {
			sendJSON(w, Response{Error: err.Error()}, http.StatusUnprocessableEntity)
			return
		}
		sendJSON(w, Response{Data: i}, http.StatusOK)
	}
}
