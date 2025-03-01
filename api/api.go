package api

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

func NewHandler() http.Handler {
	h := chi.NewMux()

	ctx := context.Background()
	clientFirestore := createClient(ctx)

	app := application{Data: make(map[string]string), FirestoreClient: clientFirestore}

	h.Get("/", HelloApi(app))
	h.Get("/api", GetData(app))
	h.Get("/api/languages", GetLanguage(app))

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

func GetLanguage(app application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		rapid_host := os.Getenv(X_RAPIDAPI_HOST)
		rapid_key := os.Getenv(X_RAPIDAPI_KEY)

		url := fmt.Sprintf("https://%s/languages", rapid_host)

		req, _ := http.NewRequest("GET", url, nil)

		req.Header.Add("x-rapidapi-key", rapid_key)
		req.Header.Add("x-rapidapi-host", rapid_host)

		res, _ := http.DefaultClient.Do(req)

		defer res.Body.Close()
		body, _ := io.ReadAll(res.Body)

		sendJSON(w, Response{Data: string(body)}, http.StatusOK)
	}
}
