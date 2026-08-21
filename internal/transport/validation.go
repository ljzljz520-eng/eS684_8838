package transport

import (
	"net/http"
	"strings"
)

func methodAllowed(r *http.Request, methods ...string) bool {
	for _, method := range methods {
		if r.Method == method {
			return true
		}
	}
	return false
}

func pathID(path, prefix string) string { return strings.TrimPrefix(strings.TrimSpace(path), prefix) }

func setCommonHeaders(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
}

func statusForError(message string) int {
	if strings.Contains(message, "not found") {
		return http.StatusNotFound
	}
	if strings.Contains(message, "required") {
		return http.StatusBadRequest
	}
	return http.StatusConflict
}
