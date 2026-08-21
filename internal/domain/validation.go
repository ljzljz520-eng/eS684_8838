package domain

import (
	"fmt"
	"strings"
)

func ValidateVisitor(record VisitorRecord) []string {
	issues := make([]string, 0, 4)
	if record.Name == "" {
		issues = append(issues, "name is required")
	}
	if record.Company == "" {
		issues = append(issues, "company is required")
	}
	if record.Host == "" {
		issues = append(issues, "host is required")
	}
	if record.VisitDate == "" {
		issues = append(issues, "visit date is required")
	}
	if len(record.Name) > 80 {
		issues = append(issues, "name is too long")
	}
	return issues
}

func ValidateBatch(batch ImportBatch) error {
	if strings.TrimSpace(batch.ID) == "" {
		return fmt.Errorf("batch id is required")
	}
	if batch.Total < 0 || batch.Valid < 0 || batch.Invalid < 0 {
		return fmt.Errorf("batch counts cannot be negative")
	}
	if batch.Valid+batch.Invalid > batch.Total {
		return fmt.Errorf("batch counts exceed total")
	}
	if !ValidBatchState(batch.State) {
		return fmt.Errorf("invalid batch state %q", batch.State)
	}
	return nil
}

func NormalizeTags(tags []string) []string {
	seen := map[string]bool{}
	out := make([]string, 0, len(tags))
	for _, tag := range tags {
		clean := strings.TrimSpace(strings.ToLower(tag))
		if clean != "" && !seen[clean] {
			seen[clean] = true
			out = append(out, clean)
		}
	}
	return out
}
