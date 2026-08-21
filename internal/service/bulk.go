package service

import (
	"errors"
	"fmt"

	"parkvisitor/internal/domain"
)

type BulkResult struct {
	Created int
	Updated int
	Failed  int
	Errors  []string
}

func (a *App) UpsertVisitors(batchID string, inputs []VisitorInput) (BulkResult, error) {
	if len(inputs) == 0 {
		return BulkResult{}, errors.New("inputs are empty")
	}
	if _, err := a.ensureBatch(batchID); err != nil {
		return BulkResult{}, err
	}
	result := BulkResult{Errors: []string{}}
	for index, input := range inputs {
		id := input.ID
		if id == "" {
			id = a.nextID("visitor")
		}
		existing, err := a.store.GetVisitor(id)
		if err != nil {
			_, createErr := a.ImportSingle(batchID, VisitorInput{ID: id, Name: input.Name, Company: input.Company, Host: input.Host, VisitDate: input.VisitDate, Notes: input.Notes, Tags: input.Tags})
			if createErr != nil {
				result.Failed++
				result.Errors = append(result.Errors, fmt.Sprintf("row %d: %v", index+1, createErr))
			} else {
				result.Created++
			}
			continue
		}
		if existing.BatchID != batchID {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("row %d belongs to another batch", index+1))
			continue
		}
		if _, updateErr := a.UpdateVisitor(id, "bulk", input); updateErr != nil {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("row %d: %v", index+1, updateErr))
		} else {
			result.Updated++
		}
	}
	return result, nil
}

func (a *App) RevalidateAll(batchID string) (domain.ImportBatch, error) {
	return a.ValidateBatch(batchID)
}
