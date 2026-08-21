package service

import "testing"

func TestBusiness37Regression(t *testing.T) {
	app := newTestApp(t)
	input := VisitorInput{ID: "rb684-37-visitor", Name: "Boundary Guest", Company: "Acme", Host: "Gate", VisitDate: "2026-08-21T00:30:00+08:00"}
	first, err := app.ImportAndValidate("RB684-37", []VisitorInput{input})
	if err != nil {
		t.Fatal(err)
	}
	if len(first.Issues) != 0 {
		t.Fatalf("first import issues=%v", first.Issues)
	}
	second, err := app.ImportAndValidate("RB684-37", []VisitorInput{input})
	if err != nil {
		t.Fatal(err)
	}
	if second.Batch.Valid != 1 || second.Records[0].VisitDate != "2026-08-21" {
		t.Fatalf("inconsistent repeated result=%+v", second)
	}
}
