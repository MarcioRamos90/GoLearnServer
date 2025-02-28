package api

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"strconv"

	"cloud.google.com/go/firestore"
)

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

func createClient(ctx context.Context) *firestore.Client {
	client, err := firestore.NewClient(ctx, firestore.DetectProjectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	return client
}

func extracIntegerQueryParam(w http.ResponseWriter, r *http.Request, key string) int {
	queryParamStr := r.URL.Query().Get(key)
	if queryParamStr == "" {
		queryParamStr = "0" // default value
	}
	queryParam, err := strconv.ParseInt(queryParamStr, 10, 64)

	if err != nil {
		slog.Error("Error to parse param", "error", err)
		sendJSON(w, Response{Error: "the param needs to be numeric"}, http.StatusBadRequest)
	}
	return int(queryParam)
}
