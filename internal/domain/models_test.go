package domain

import "testing"

func TestVisitorNormalization(t *testing.T) {
	record := NewVisitorRecord("v1", "b1", "  Lin   Mei ", "Acme", "Host", "2026-08-21", "now")
	if record.Name != "Lin Mei" {
		t.Fatalf("name=%q", record.Name)
	}
	record = record.AddTag("  vip ")
	record = record.AddTag("vip")
	if len(record.Tags) != 1 || !record.HasTag("vip") {
		t.Fatalf("tags=%v", record.Tags)
	}
}

func TestValidationRules(t *testing.T) {
	record := VisitorRecord{ID: "v", BatchID: "b"}
	issues := ValidateVisitor(record)
	if len(issues) != 4 {
		t.Fatalf("issues=%v", issues)
	}
	if ValidRecordStatus("unknown") {
		t.Fatal("unknown status accepted")
	}
	if !CanPublish(BatchConfirmed) {
		t.Fatal("confirmed batch not publishable")
	}
}
