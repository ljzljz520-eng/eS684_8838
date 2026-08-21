package service

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"parkvisitor/internal/domain"
)

type BatchOverview struct {
	Batch         domain.ImportBatch `json:"batch"`
	VisitorCount  int                `json:"visitor_count"`
	PendingReview int                `json:"pending_review"`
	PendingTasks  int                `json:"pending_tasks"`
	LastAction    string             `json:"last_action"`
}

func (a *App) Overview(batchID string) (BatchOverview, error) {
	batch, err := a.store.GetBatch(batchID)
	if err != nil {
		return BatchOverview{}, err
	}
	records, err := a.store.ListVisitors(batchID)
	if err != nil {
		return BatchOverview{}, err
	}
	tasks, err := a.store.ListTasks(batchID)
	if err != nil {
		return BatchOverview{}, err
	}
	events, err := a.store.ListAudit(batchID)
	if err != nil {
		return BatchOverview{}, err
	}
	overview := BatchOverview{Batch: batch, VisitorCount: len(records)}
	for _, record := range records {
		if record.Status == domain.StatusNeedsReview {
			overview.PendingReview++
		}
	}
	for _, task := range tasks {
		if task.State == domain.TaskOpen {
			overview.PendingTasks++
		}
	}
	if len(events) > 0 {
		overview.LastAction = events[len(events)-1].Action
	}
	return overview, nil
}

func (a *App) FindDuplicates(batchID string) ([][]domain.VisitorRecord, error) {
	records, err := a.store.ListVisitors(batchID)
	if err != nil {
		return nil, err
	}
	groups := map[string][]domain.VisitorRecord{}
	for _, record := range records {
		key := strings.ToLower(strings.TrimSpace(record.Name + "|" + record.Company))
		groups[key] = append(groups[key], record)
	}
	result := [][]domain.VisitorRecord{}
	for _, group := range groups {
		if len(group) > 1 {
			sort.Slice(group, func(i, j int) bool { return group[i].ID < group[j].ID })
			result = append(result, group)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i][0].ID < result[j][0].ID })
	return result, nil
}

func (a *App) AddTag(recordID, tag, actor string) (domain.VisitorRecord, error) {
	record, err := a.store.GetVisitor(recordID)
	if err != nil {
		return record, err
	}
	if strings.TrimSpace(tag) == "" {
		return record, errors.New("tag is required")
	}
	record = record.AddTag(tag)
	record.UpdatedAt = a.clock.NowString()
	if err := a.store.SaveVisitor(record); err != nil {
		return record, err
	}
	_ = a.audit(record.BatchID, record.ID, "tag_added", actor, tag)
	return record, nil
}

func (a *App) RemoveTag(recordID, tag, actor string) (domain.VisitorRecord, error) {
	record, err := a.store.GetVisitor(recordID)
	if err != nil {
		return record, err
	}
	clean := strings.ToLower(strings.TrimSpace(tag))
	filtered := make([]string, 0, len(record.Tags))
	for _, item := range record.Tags {
		if strings.ToLower(item) != clean {
			filtered = append(filtered, item)
		}
	}
	if len(filtered) == len(record.Tags) {
		return record, fmt.Errorf("tag %s not found", tag)
	}
	record.Tags = filtered
	record.UpdatedAt = a.clock.NowString()
	if err := a.store.SaveVisitor(record); err != nil {
		return record, err
	}
	_ = a.audit(record.BatchID, record.ID, "tag_removed", actor, tag)
	return record, nil
}

func (a *App) RenameBatch(batchID, reference string) (domain.ImportBatch, error) {
	batch, err := a.store.GetBatch(batchID)
	if err != nil {
		return batch, err
	}
	reference = strings.TrimSpace(reference)
	if reference == "" {
		return batch, errors.New("reference is required")
	}
	batch.Reference = reference
	if err := a.store.SaveBatch(batch); err != nil {
		return batch, err
	}
	_ = a.audit(batchID, "", "batch_renamed", "system", reference)
	return batch, nil
}
