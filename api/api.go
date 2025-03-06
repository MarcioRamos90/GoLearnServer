package api

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"server/types"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func NewHandler() http.Handler {

	ctx := context.Background()
	clientFirestore := createClient(ctx)
	app := types.Application{Envs: types.Envs{}, FirestoreClient: clientFirestore}

	api := chi.NewMux()
	api.Use(middleware.Recoverer)
	api.Use(middleware.RequestID)
	api.Use(middleware.Logger)
	api.Get("/", HelloApi(app))
	api.Route("/api", func(r chi.Router) {
		app.Envs.Rapid_host = os.Getenv(X_RAPIDAPI_HOST)
		app.Envs.Rapid_key = os.Getenv(X_RAPIDAPI_KEY)
		r.Get("/languages", GetLanguage(app))
		r.Post("/submit", PostSubmitCode(app))
		r.Get("/submit/{submitionId}", GetSubmitionCode(app))
		r.Post("/wishlist", PostWishList(app))
	})

	return api
}

func HelloApi(app types.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		SendJSON(w, types.Response{Data: "Hello API!"}, http.StatusOK)
	}
}

func PostWishList(app types.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var wishListSubscription types.WishListSubscription
		if err := json.NewDecoder(r.Body).Decode(&wishListSubscription); err != nil {
			slog.Error("error on body decode", "error", err)
			SendJSON(w, types.Response{Error: "something went wrong"}, http.StatusInternalServerError)
			return
		}

		if wishListSubscription.Email == "" {
			SendJSON(w, types.Response{Error: "please send email"}, http.StatusBadRequest)
			return
		}

		subsRef := app.FirestoreClient.Collection("subscriptions")

		// verify if email already exists
		query := subsRef.Select().Where("Email", "==", wishListSubscription.Email)
		doc, _ := query.Documents(r.Context()).Next()

		if doc != nil && doc.Exists() {
			SendJSON(w, types.Response{Error: "this email already exists on our registration"}, http.StatusOK)
			return
		}

		// send to firestore
		_, _, err := subsRef.Add(r.Context(), wishListSubscription)

		if err != nil {
			slog.Error("error on saving email", "error", err)
			SendJSON(w, types.Response{Error: "error on saving email"}, http.StatusBadRequest)
			return
		}

		SendJSON(w, types.Response{Data: "subscribed"}, http.StatusCreated)
	}
}
