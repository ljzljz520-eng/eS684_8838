package policy

import (
	"fmt"
	"strings"

	"parkvisitor/internal/domain"
)

type Decision struct {
	Allowed bool
	Reasons []string
}

func (r Rules) Decide(record domain.VisitorRecord) Decision {
	reasons := []string{}
	if err := r.ValidateInput(record); err != nil {
		reasons = append(reasons, err.Error())
	}
	if record.Status == domain.StatusRejected {
		reasons = append(reasons, "rejected records need a new review")
	}
	return Decision{Allowed: len(reasons) == 0, Reasons: reasons}
}

func Explain(decision Decision) string {
	if decision.Allowed {
		return "allowed"
	}
	return strings.Join(decision.Reasons, "; ")
}

func (r Rules) ValidateBatch(records []domain.VisitorRecord) error {
	if !r.AcceptBatch(len(records)) {
		return fmt.Errorf("batch size %d exceeds policy maximum %d", len(records), r.MaxBatchSize)
	}
	for _, record := range records {
		if err := r.ValidateInput(record); err != nil {
			return fmt.Errorf("record %s: %w", record.ID, err)
		}
	}
	return nil
}

func (r Rules) AllowedStatus(status string) bool {
	switch status {
	case domain.StatusImported, domain.StatusValidated, domain.StatusNeedsReview, domain.StatusApproved, domain.StatusPublished, domain.StatusArchived:
		return true
	case domain.StatusRejected:
		return false
	}
	return false
}
