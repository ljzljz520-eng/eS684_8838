package transport

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"parkvisitor/internal/clock"
	"parkvisitor/internal/service"
	"parkvisitor/internal/storage"
)

func TestHTTPWorkflow(t *testing.T) {
	app := serviceTestApp(t)
	handler := NewHandler(app)
	body := `{"id":"batch-http","reference":"RB-http","source":"api","inputs":[{"id":"http-1","name":"Guest","company":"Acme","host":"Host","visit_date":"2026-08-21"}]}`
	request := httptest.NewRequest(http.MethodPost, "/batches", strings.NewReader(body))
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusCreated {
		t.Fatalf("create status=%d body=%s", response.Code, response.Body.String())
	}
	request = httptest.NewRequest(http.MethodGet, "/visitors?batch=batch-http", nil)
	response = httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusOK || !strings.Contains(response.Body.String(), "http-1") {
		t.Fatalf("search status=%d body=%s", response.Code, response.Body.String())
	}
}

func serviceTestApp(t *testing.T) *service.App { t.Helper(); return serviceTestAppWithCleanup(t) }

func serviceTestAppWithCleanup(t *testing.T) *service.App {
	t.Helper()
	store, err := storage.Open(t.TempDir() + "/http.db")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { store.Close() })
	app, err := service.NewApp(store, clock.Fixed())
	if err != nil {
		t.Fatal(err)
	}
	return app
}
