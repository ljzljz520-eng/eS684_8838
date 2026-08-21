package service

import (
	"testing"

	"parkvisitor/internal/clock"
	"parkvisitor/internal/storage"
)

func newTestApp(t *testing.T) *App {
	t.Helper()
	store, err := storage.Open(t.TempDir() + "/data.db")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { store.Close() })
	app, err := NewApp(store, clock.Fixed())
	if err != nil {
		t.Fatal(err)
	}
	return app
}

func TestWorkflowCreateReviewArchive(t *testing.T) {
	app := newTestApp(t)
	input := VisitorInput{ID: "visitor-a", Name: "Lin Mei", Company: "Acme", Host: "Zhou", VisitDate: "2026-08-21", Tags: []string{"vip"}}
	result, err := app.ImportAndValidate("batch-a", []VisitorInput{input})
	if err != nil {
		t.Fatal(err)
	}
	if result.Batch.State != "confirmed" {
		t.Fatalf("state=%s", result.Batch.State)
	}
	if _, err := app.ReviewBatch("batch-a", "reviewer"); err != nil {
		t.Fatal(err)
	}
	batch, err := app.PublishBatch("batch-a", "publisher")
	if err != nil {
		t.Fatal(err)
	}
	if batch.State != "published" {
		t.Fatalf("published state=%s", batch.State)
	}
	batch, err = app.ArchiveBatch("batch-a", "archiver")
	if err != nil {
		t.Fatal(err)
	}
	if batch.State != "archived" {
		t.Fatalf("archived state=%s", batch.State)
	}
}
