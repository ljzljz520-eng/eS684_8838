package transport

import (
	"fmt"
	"net/http"
	"strings"
)

type queryOptions struct {
	Text    string
	Company string
	Status  string
	Tag     string
	Batch   string
}

func parseQuery(r *http.Request) queryOptions {
	values := r.URL.Query()
	return queryOptions{Text: strings.TrimSpace(values.Get("q")), Company: strings.TrimSpace(values.Get("company")), Status: strings.TrimSpace(values.Get("status")), Tag: strings.TrimSpace(values.Get("tag")), Batch: strings.TrimSpace(values.Get("batch"))}
}

func acceptsJSON(r *http.Request) bool {
	return strings.Contains(r.Header.Get("Accept"), "application/json") || r.Header.Get("Accept") == ""
}

func requestID(r *http.Request) string {
	value := strings.TrimSpace(r.Header.Get("X-Request-ID"))
	if value == "" {
		return "local-request"
	}
	return value
}

func pagination(r *http.Request) (int, int) {
	limit, offset := 50, 0
	if value := r.URL.Query().Get("limit"); value != "" {
		if _, err := fmt.Sscanf(value, "%d", &limit); err != nil || limit < 1 {
			limit = 50
		}
	}
	if value := r.URL.Query().Get("offset"); value != "" {
		if _, err := fmt.Sscanf(value, "%d", &offset); err != nil || offset < 0 {
			offset = 0
		}
	}
	if limit > 200 {
		limit = 200
	}
	return limit, offset
}

func slicePage[T any](items []T, limit, offset int) []T {
	if offset >= len(items) {
		return []T{}
	}
	end := offset + limit
	if end > len(items) {
		end = len(items)
	}
	return items[offset:end]
}
