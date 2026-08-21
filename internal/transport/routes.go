package transport

import (
	"encoding/json"
	"net/http"
	"strings"

	"parkvisitor/internal/service"
)

func registerRoutes(mux *http.ServeMux, app *service.App) {
	RegisterHealth(mux)
	mux.HandleFunc("/batches", batchHandler(app))
	mux.HandleFunc("/visitors", visitorHandler(app))
	mux.HandleFunc("/reports/", reportHandler(app))
}

func batchHandler(app *service.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		var payload struct {
			ID        string                 `json:"id"`
			Reference string                 `json:"reference"`
			Source    string                 `json:"source"`
			Inputs    []service.VisitorInput `json:"inputs"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		if payload.ID == "" {
			writeError(w, http.StatusBadRequest, "id is required")
			return
		}
		if len(payload.Inputs) == 0 {
			batch, err := app.CreateBatch(payload.ID, payload.Reference, payload.Source)
			if err != nil {
				writeError(w, http.StatusBadRequest, err.Error())
				return
			}
			writeJSON(w, http.StatusCreated, batch)
			return
		}
		result, err := app.ImportAndValidate(payload.ID, payload.Inputs)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeJSON(w, http.StatusCreated, result)
	}
}

func visitorHandler(app *service.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		query := service.Query{Text: r.URL.Query().Get("q"), Company: r.URL.Query().Get("company"), Status: r.URL.Query().Get("status"), Tag: r.URL.Query().Get("tag"), BatchID: r.URL.Query().Get("batch")}
		records, err := app.SearchVisitors(query)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, records)
	}
}

func reportHandler(app *service.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		batchID := strings.TrimPrefix(r.URL.Path, "/reports/")
		if batchID == "" {
			writeError(w, http.StatusBadRequest, "batch id is required")
			return
		}
		result, err := app.GenerateReport(batchID)
		if err != nil {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, result)
	}
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
