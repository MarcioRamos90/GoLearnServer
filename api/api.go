package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
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

	h.Post("/api/wishlist", PostWishList(app))

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

func PostWishList(app application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var wishListSubscription WishListSubscription
		if err := json.NewDecoder(r.Body).Decode(&wishListSubscription); err != nil {
			slog.Error("error on body decode", "error", err)
			sendJSON(w, Response{Error: "something went wrong"}, http.StatusInternalServerError)
			return
		}

		if wishListSubscription.Email == "" {
			sendJSON(w, Response{Error: "please send email"}, http.StatusBadRequest)
			return
		}

		subsRef := app.FirestoreClient.Collection("subscriptions")

		// verify if email already exists
		query := subsRef.Select().Where("Email", "==", wishListSubscription.Email)
		doc, _ := query.Documents(r.Context()).Next()

		if doc != nil && doc.Exists() {
			sendJSON(w, Response{Error: "this email already exists on our registration"}, http.StatusOK)
			return
		}

		// send to firestore
		_, _, err := subsRef.Add(r.Context(), wishListSubscription)

		if err != nil {
			slog.Error("error on saving email", "error", err)
			sendJSON(w, Response{Error: "error on saving email"}, http.StatusBadRequest)
			return
		}

		sendJSON(w, Response{Data: "subscribed"}, http.StatusCreated)
	}
}
