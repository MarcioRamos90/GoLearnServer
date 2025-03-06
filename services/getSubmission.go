package services

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"server/types"
)

func NewRequest(app types.Application, method string, url string, payload io.Reader) *http.Request {
	req, _ := http.NewRequest(method, url, payload)

	req.Header.Add("x-rapidapi-key", app.Envs.Rapid_key)
	req.Header.Add("x-rapidapi-host", app.Envs.Rapid_host)

	if payload != nil {
		req.Header.Add("Content-Type", "application/json")
	}

	return req
}

func GetSubmissionById(app types.Application, id string) ([]byte, error) {
	url := fmt.Sprintf("https://%s/submissions/%s?base64_encoded=false&fields=stdout", app.Envs.Rapid_host, id)
	slog.Info("GetSubmissionById", "url", url)
	res, _ := http.DefaultClient.Do(NewRequest(app, http.MethodGet, url, nil))
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	slog.Info("GetSubmissionById", "body", body)

	return body, nil
}
