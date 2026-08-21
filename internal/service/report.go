package service

import (
	"sort"

	"parkvisitor/internal/domain"
)

func (a *App) GenerateReport(batchID string) (domain.Report, error) {
	batch, err := a.store.GetBatch(batchID)
	if err != nil {
		return domain.Report{}, err
	}
	records, err := a.store.ListVisitors(batchID)
	if err != nil {
		return domain.Report{}, err
	}
	tasks, err := a.store.ListTasks(batchID)
	if err != nil {
		return domain.Report{}, err
	}
	audit, err := a.store.ListAudit(batchID)
	if err != nil {
		return domain.Report{}, err
	}
	report := domain.Report{BatchID: batchID, State: batch.State, Total: len(records), ByStatus: map[string]int{}, ByCompany: map[string]int{}, AuditCount: len(audit)}
	for _, record := range records {
		report.ByStatus[record.Status]++
		report.ByCompany[record.Company]++
	}
	for _, task := range tasks {
		if task.State == domain.TaskOpen {
			report.PendingTasks++
		}
	}
	return report, nil
}

func (a *App) ExportVisitors(batchID string) ([]domain.VisitorRecord, error) {
	records, err := a.store.ListVisitors(batchID)
	if err != nil {
		return nil, err
	}
	sort.Slice(records, func(i, j int) bool { return records[i].ID < records[j].ID })
	return records, nil
}
