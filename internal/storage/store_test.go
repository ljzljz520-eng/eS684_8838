package storage

import (
	"testing"

	"parkvisitor/internal/domain"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	path := t.TempDir() + "/visitors.db"
	store, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	batch := domain.ImportBatch{ID: "batch-1", Reference: "RB684-1", Source: "test", State: domain.BatchConfirmed, Total: 1, Valid: 1}
	record := domain.NewVisitorRecord("record-1", batch.ID, "Lin Mei", "Acme", "Host", "2026-08-21", "2026-08-21T08:00:00Z")
	event := domain.AuditEvent{ID: "audit-1", BatchID: batch.ID, RecordID: record.ID, Action: "import", Actor: "test", At: record.CreatedAt}
	attachment := domain.Attachment{ID: "attachment-1", RecordID: record.ID, Name: "badge.pdf", Digest: "abc", Size: 3}
	task := domain.CollaborationTask{ID: "task-1", BatchID: batch.ID, RecordID: record.ID, Assignee: "ops", State: domain.TaskOpen}
	for _, save := range []func() error{func() error { return store.SaveBatch(batch) }, func() error { return store.SaveVisitor(record) }, func() error { return store.SaveAudit(event) }, func() error { return store.SaveAttachment(attachment) }, func() error { return store.SaveTask(task) }} {
		if err := save(); err != nil {
			t.Fatal(err)
		}
	}
	if err := store.Close(); err != nil {
		t.Fatal(err)
	}
	reopened, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer reopened.Close()
	if got, err := reopened.GetBatch(batch.ID); err != nil || got.Reference != batch.Reference {
		t.Fatalf("batch=%+v err=%v", got, err)
	}
	if got, err := reopened.GetVisitor(record.ID); err != nil || got.Company != record.Company {
		t.Fatalf("record=%+v err=%v", got, err)
	}
	if got, err := reopened.ListAudit(batch.ID); err != nil || len(got) != 1 {
		t.Fatalf("audit=%v err=%v", got, err)
	}
	if got, err := reopened.ListAttachments(record.ID); err != nil || len(got) != 1 {
		t.Fatalf("attachments=%v err=%v", got, err)
	}
	if got, err := reopened.GetTask(task.ID); err != nil || got.Assignee != task.Assignee {
		t.Fatalf("task=%+v err=%v", got, err)
	}
}
