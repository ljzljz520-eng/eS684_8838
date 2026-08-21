package service

import (
	"fmt"
	"strings"

	"parkvisitor/internal/domain"
	"parkvisitor/internal/policy"
)

func (a *App) ValidateRecord(record domain.VisitorRecord) error {
	issues := domain.ValidateVisitor(record)
	if len(issues) > 0 {
		return fmt.Errorf("record %s: %s", record.ID, strings.Join(issues, ", "))
	}
	decision := policy.DefaultRules().Decide(record)
	if !decision.Allowed {
		return fmt.Errorf("record %s: %s", record.ID, policy.Explain(decision))
	}
	return nil
}

func (a *App) ValidateRecords(records []domain.VisitorRecord) map[string][]string {
	result := map[string][]string{}
	for _, record := range records {
		if err := a.ValidateRecord(record); err != nil {
			result[record.ID] = []string{err.Error()}
		}
	}
	return result
}

func (a *App) RebuildBatchCounts(batchID string) (domain.ImportBatch, error) {
	batch, err := a.store.GetBatch(batchID)
	if err != nil {
		return batch, err
	}
	records, err := a.store.ListVisitors(batchID)
	if err != nil {
		return batch, err
	}
	batch.Total = len(records)
	batch.Valid = 0
	batch.Invalid = 0
	for _, record := range records {
		if a.ValidateRecord(record) == nil {
			batch.Valid++
		} else {
			batch.Invalid++
		}
	}
	if batch.Invalid == 0 {
		batch.State = domain.BatchConfirmed
	} else {
		batch.State = domain.BatchValidated
	}
	if err := a.store.SaveBatch(batch); err != nil {
		return batch, err
	}
	return batch, nil
}

func (a *App) CanTransition(recordID, next string) (bool, error) {
	record, err := a.store.GetVisitor(recordID)
	if err != nil {
		return false, err
	}
	return domain.AllowedTransition(record.Status, next), nil
}
