package api

import "cloud.google.com/go/firestore"

type Response struct {
	Error string `json:"error,omitempty"`
	Data  any    `json:"data,omitempty"`
}

type application struct {
	Data            map[string]string
	FirestoreClient *firestore.Client
}

type WishListSubscription struct {
	Email string `json:"email"`
}
