package service

import "testing"

func TestWorkflowSearchUpdatePublish(t *testing.T) {
	app := newTestApp(t)
	_, err := app.ImportAndValidate("batch-q", []VisitorInput{{ID: "visitor-q", Name: "Qiao Lu", Company: "North", Host: "Wang", VisitDate: "2026-08-21", Tags: []string{"partner"}}})
	if err != nil {
		t.Fatal(err)
	}
	found, err := app.SearchVisitors(Query{Text: "qiao", BatchID: "batch-q"})
	if err != nil || len(found) != 1 {
		t.Fatalf("found=%v err=%v", found, err)
	}
	updated, err := app.UpdateVisitor("visitor-q", "editor", VisitorInput{Notes: "updated"})
	if err != nil {
		t.Fatal(err)
	}
	if updated.Notes != "updated" {
		t.Fatal("note not updated")
	}
	if _, err := app.PublishBatch("batch-q", "publisher"); err != nil {
		t.Fatal(err)
	}
	found, err = app.SearchVisitors(Query{Status: "published"})
	if err != nil || len(found) != 1 {
		t.Fatalf("published=%v err=%v", found, err)
	}
}
