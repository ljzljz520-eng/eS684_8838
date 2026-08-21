package service

import (
	"errors"
	"fmt"
	"strings"

	"parkvisitor/internal/domain"
)

func (a *App) ImportAndValidate(batchID string, inputs []VisitorInput) (ImportSummary, error) {
	batch, err := a.ensureBatch(batchID)
	if err != nil {
		return ImportSummary{}, err
	}
	if len(inputs) == 0 {
		return ImportSummary{}, errors.New("at least one visitor is required")
	}
	if batch.State == domain.BatchConfirmed || batch.State == domain.BatchPublished {
		return a.summaryFromExisting(batchID, batch)
	}
	issues := map[string][]string{}
	records := make([]domain.VisitorRecord, 0, len(inputs))
	valid := 0
	for index, input := range inputs {
		id := input.ID
		if id == "" {
			id = a.nextID("visitor")
		}
		date, dateErr := a.parseBusinessDate(input.VisitDate)
		record := domain.NewVisitorRecord(id, batchID, input.Name, input.Company, input.Host, date, a.clock.NowString())
		record.Notes = strings.TrimSpace(input.Notes)
		record.Tags = domain.NormalizeTags(input.Tags)
		itemIssues := domain.ValidateVisitor(record)
		if dateErr != nil {
			itemIssues = append(itemIssues, dateErr.Error())
		}
		if len(itemIssues) > 0 {
			record.Status = domain.StatusNeedsReview
			issues[id] = itemIssues
		} else {
			record.Status = domain.StatusValidated
			valid++
		}
		if err := a.store.SaveVisitor(record); err != nil {
			return ImportSummary{}, err
		}
		records = append(records, record)
		_ = a.audit(batchID, id, "import_validated", "system", fmt.Sprintf("row=%d", index+1))
	}
	batch.Total = len(records)
	batch.Valid = valid
	batch.Invalid = len(records) - valid
	batch.BusinessDate = a.clock.Date()
	batch.State = domain.BatchValidated
	if batch.Invalid == 0 {
		batch.State = domain.BatchConfirmed
		batch.ConfirmedAt = a.clock.NowString()
	}
	if err := a.store.SaveBatch(batch); err != nil {
		return ImportSummary{}, err
	}
	return ImportSummary{Batch: batch, Records: records, Issues: issues}, nil
}

func (a *App) summaryFromExisting(batchID string, batch domain.ImportBatch) (ImportSummary, error) {
	records, err := a.store.ListVisitors(batchID)
	if err != nil {
		return ImportSummary{}, err
	}
	issues := map[string][]string{}
	for _, record := range records {
		if record.Status == domain.StatusNeedsReview {
			issues[record.ID] = []string{"record requires review"}
		}
	}
	return ImportSummary{Batch: batch, Records: records, Issues: issues}, nil
}

func (a *App) ValidateBatch(batchID string) (domain.ImportBatch, error) {
	batch, err := a.store.GetBatch(batchID)
	if err != nil {
		return batch, err
	}
	records, err := a.store.ListVisitors(batchID)
	if err != nil {
		return batch, err
	}
	valid := 0
	for _, record := range records {
		if len(domain.ValidateVisitor(record)) == 0 {
			valid++
		}
	}
	batch.Total = len(records)
	batch.Valid = valid
	batch.Invalid = len(records) - valid
	if batch.Invalid == 0 {
		batch.State = domain.BatchConfirmed
		batch.ConfirmedAt = a.clock.NowString()
	} else {
		batch.State = domain.BatchValidated
	}
	if err := a.store.SaveBatch(batch); err != nil {
		return batch, err
	}
	_ = a.audit(batchID, "", "batch_validated", "system", fmt.Sprintf("valid=%d", valid))
	return batch, nil
}

func (a *App) ImportSingle(batchID string, input VisitorInput) (domain.VisitorRecord, error) {
	result, err := a.ImportAndValidate(batchID, []VisitorInput{input})
	if err != nil {
		return domain.VisitorRecord{}, err
	}
	if len(result.Records) != 1 {
		return domain.VisitorRecord{}, errors.New("single import returned unexpected count")
	}
	return result.Records[0], nil
}
