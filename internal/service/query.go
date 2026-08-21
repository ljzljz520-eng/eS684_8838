package service

import (
	"sort"
	"strings"

	"parkvisitor/internal/domain"
)

type Query struct {
	Text    string
	Company string
	Status  string
	Tag     string
	BatchID string
}

func (a *App) SearchVisitors(query Query) ([]domain.VisitorRecord, error) {
	records, err := a.store.ListVisitors(query.BatchID)
	if err != nil {
		return nil, err
	}
	result := make([]domain.VisitorRecord, 0, len(records))
	needle := strings.ToLower(strings.TrimSpace(query.Text))
	for _, record := range records {
		if needle != "" && !strings.Contains(strings.ToLower(record.Name+" "+record.Company+" "+record.Host), needle) {
			continue
		}
		if query.Company != "" && !strings.EqualFold(record.Company, query.Company) {
			continue
		}
		if query.Status != "" && record.Status != query.Status {
			continue
		}
		if query.Tag != "" && !record.HasTag(query.Tag) {
			continue
		}
		result = append(result, record)
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].VisitDate == result[j].VisitDate {
			return result[i].Name < result[j].Name
		}
		return result[i].VisitDate < result[j].VisitDate
	})
	return result, nil
}

func (a *App) GetVisitor(id string) (domain.VisitorRecord, error) { return a.store.GetVisitor(id) }

func (a *App) UpdateVisitor(id, actor string, update VisitorInput) (domain.VisitorRecord, error) {
	record, err := a.store.GetVisitor(id)
	if err != nil {
		return record, err
	}
	if update.Name != "" {
		record.Name = domain.NormalizeName(update.Name)
	}
	if update.Company != "" {
		record.Company = strings.TrimSpace(update.Company)
	}
	if update.Host != "" {
		record.Host = strings.TrimSpace(update.Host)
	}
	if update.VisitDate != "" {
		date, parseErr := a.parseBusinessDate(update.VisitDate)
		if parseErr != nil {
			return record, parseErr
		}
		record.VisitDate = date
	}
	if update.Notes != "" {
		record.Notes = strings.TrimSpace(update.Notes)
	}
	if update.Tags != nil {
		record.Tags = domain.NormalizeTags(update.Tags)
	}
	if len(domain.ValidateVisitor(record)) > 0 {
		return record, domainError(record)
	}
	record.UpdatedAt = a.clock.NowString()
	if err := a.store.SaveVisitor(record); err != nil {
		return record, err
	}
	_ = a.audit(record.BatchID, id, "updated", actor, "visitor details updated")
	return record, nil
}

func domainError(record domain.VisitorRecord) error { return &validationError{recordID: record.ID} }

type validationError struct{ recordID string }

func (e *validationError) Error() string { return "visitor record remains incomplete: " + e.recordID }
