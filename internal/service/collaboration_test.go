package service

import "testing"

func TestCollaborationAndAttachment(t *testing.T) {
	app := newTestApp(t)
	_, err := app.ImportAndValidate("batch-c", []VisitorInput{{ID: "c1", Name: "Chen", Company: "Acme", Host: "Host", VisitDate: "2026-08-21"}})
	if err != nil {
		t.Fatal(err)
	}
	attachment, err := app.AttachDocument("c1", "badge.pdf", "badge", "pdf")
	if err != nil || attachment.Size != 3 {
		t.Fatalf("attachment=%+v err=%v", attachment, err)
	}
	task, err := app.AddCollaboration("batch-c", "c1", "security", "check", "2026-08-22")
	if err != nil {
		t.Fatal(err)
	}
	tasks, err := app.ListCollaboration("batch-c")
	if err != nil || len(tasks) != 1 {
		t.Fatalf("tasks=%v err=%v", tasks, err)
	}
	if _, err := app.CompleteTask(task.ID, "security"); err != nil {
		t.Fatal(err)
	}
}
