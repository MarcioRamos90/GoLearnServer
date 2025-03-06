package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"server/services"
	"server/types"
	"strings"

	"github.com/go-chi/chi/v5"
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

type CodeSubmission struct {
	Code  string `json:"code"`
	Title string `json:"title,omitempty"`
	ID    string `json:"id,omitempty"`
}

type RequestBody struct {
	LanguageID int    `json:"language_id"`
	SourceCode string `json:"source_code"`
	Stdin      string `json:"stdin,omitempty"`
}

type RapidPostResponseApi struct {
	Token string `json:"token"`
}

type RapidGetResponseApi struct {
	Stdout string `json:"stdout"`
}

func PostSubmitCode(app types.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var codeSubmission CodeSubmission
		if err := json.NewDecoder(r.Body).Decode(&codeSubmission); err != nil {
			slog.Error("error on body decode", "error", err)
			SendJSON(w, types.Response{Error: "something went wrong"}, http.StatusInternalServerError)
			return
		}

		url := fmt.Sprintf("https://%s/submissions?base64_encoded=true&wait=false&fields=*", app.Envs.Rapid_host)
		bodyData := RequestBody{LanguageID: 106, SourceCode: codeSubmission.Code, Stdin: "SnVkZ2Uw"}
		bodyByte, _ := json.Marshal(bodyData)
		bodyReader := strings.NewReader(string(bodyByte))
		slog.Info("Request data", "url", url, "body", bodyReader)

		req := NewRequest(app, http.MethodPost, url, bodyReader)

		res, _ := http.DefaultClient.Do(req)
		defer res.Body.Close()
		body, _ := io.ReadAll(res.Body)
		slog.Info("body response", "body", body)

		var mapPostData RapidPostResponseApi
		reader := strings.NewReader(string(body))

		if err := json.NewDecoder(reader).Decode(&mapPostData); err != nil {
			slog.Error("error on body decode", "error", err)
			SendJSON(w, types.Response{Error: "something went wrong"}, http.StatusInternalServerError)
			return
		}
		slog.Info("marshaled data", "data", mapPostData)

		SendJSON(w, types.Response{Data: mapPostData.Token}, http.StatusOK)
	}
}

func GetSubmitionCode(app types.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		submitionId := chi.URLParam(r, "submitionId")

		body, _ := services.GetSubmissionById(app, submitionId)
		var mapGetData RapidGetResponseApi
		reader := strings.NewReader(string(body))
		if err := json.NewDecoder(reader).Decode(&mapGetData); err != nil {
			slog.Error("error on body decode", "error", err)
			SendJSON(w, types.Response{Error: "something went wrong"}, http.StatusInternalServerError)
			return
		}

		SendJSON(w, types.Response{Data: mapGetData.Stdout}, http.StatusOK)
	}
}

func GetLanguage(app types.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		url := fmt.Sprintf("https://%s/languages", app.Envs.Rapid_host)

		req := NewRequest(app, http.MethodGet, url, nil)
		res, _ := http.DefaultClient.Do(req)
		slog.Info("Request data", "url", url, "key", app.Envs.Rapid_key)
		defer res.Body.Close()
		body, _ := io.ReadAll(res.Body)

		SendJSON(w, types.Response{Data: string(body)}, http.StatusOK)
	}
}
