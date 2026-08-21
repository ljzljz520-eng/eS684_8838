package transport

import (
	"net/http"

	"parkvisitor/internal/service"
)

func NewHandler(app *service.App) http.Handler {
	mux := http.NewServeMux()
	registerRoutes(mux, app)
	return requestLog(mux)
}

func requestLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Visitor-Service", "park-sync")
		next.ServeHTTP(w, r)
	})
}

func RegisterHealth(mux *http.ServeMux) {
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusNoContent) })
}
