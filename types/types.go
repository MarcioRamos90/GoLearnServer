package types

import "cloud.google.com/go/firestore"

type Response struct {
	Error string `json:"error,omitempty"`
	Data  any    `json:"data,omitempty"`
}

type Envs struct {
	Rapid_host string
	Rapid_key  string
}

type Application struct {
	FirestoreClient *firestore.Client
	Envs            Envs
}

type WishListSubscription struct {
	Email string `json:"email"`
}
