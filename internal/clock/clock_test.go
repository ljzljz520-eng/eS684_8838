package clock

import (
	"testing"
	"time"
)

func TestFixedClockIsDeterministic(t *testing.T) {
	first := Fixed()
	second := Fixed()
	if !first.Now().Equal(second.Now()) {
		t.Fatal("fixed clocks differ")
	}
	if first.Date() != "2026-08-21" {
		t.Fatalf("date=%s", first.Date())
	}
	if !first.AtDate("2026-08-21").Equal(time.Date(2026, 8, 21, 0, 0, 0, 0, time.UTC)) {
		t.Fatal("date conversion mismatch")
	}
}
