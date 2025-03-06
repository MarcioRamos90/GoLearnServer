package api

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"server/types"
	"strconv"

	"cloud.google.com/go/firestore"
)

func SendJSON(w http.ResponseWriter, resp types.Response, status int) {
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

func extracIntegerFromQueryParam(w http.ResponseWriter, r *http.Request, key string) int {
	queryParamStr := r.URL.Query().Get(key)
	if queryParamStr == "" {
		queryParamStr = "0"
	}
	queryParam, err := strconv.ParseInt(queryParamStr, 10, 64)

	if err != nil {
		slog.Error("Error to parse param", "error", err)
		SendJSON(w, types.Response{Error: "the param needs to be numeric"}, http.StatusBadRequest)
	}
	return int(queryParam)
}
