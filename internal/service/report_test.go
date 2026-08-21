package service

import "testing"

func TestWorkflowImportReport(t *testing.T) {
	app := newTestApp(t)
	_, err := app.ImportAndValidate("batch-r", []VisitorInput{{ID: "r1", Name: "A", Company: "Acme", Host: "H", VisitDate: "2026-08-21"}, {ID: "r2", Name: "B", Company: "Beta", Host: "H", VisitDate: "2026-08-21"}})
	if err != nil {
		t.Fatal(err)
	}
	task, err := app.AddCollaboration("batch-r", "r1", "ops", "verify badge", "2026-08-22")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := app.CompleteTask(task.ID, "ops"); err != nil {
		t.Fatal(err)
	}
	result, err := app.GenerateReport("batch-r")
	if err != nil {
		t.Fatal(err)
	}
	if result.Total != 2 || result.ByCompany["Acme"] != 1 || result.PendingTasks != 0 {
		t.Fatalf("report=%+v", result)
	}
}
