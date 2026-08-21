package service

import (
	"sort"
	"strings"

	"parkvisitor/internal/domain"
)

type AuditFilter struct {
	Actor    string
	Action   string
	RecordID string
}

func FilterAudit(events []domain.AuditEvent, filter AuditFilter) []domain.AuditEvent {
	result := make([]domain.AuditEvent, 0, len(events))
	for _, event := range events {
		if filter.Actor != "" && !strings.EqualFold(event.Actor, filter.Actor) {
			continue
		}
		if filter.Action != "" && event.Action != filter.Action {
			continue
		}
		if filter.RecordID != "" && event.RecordID != filter.RecordID {
			continue
		}
		result = append(result, event)
	}
	sort.SliceStable(result, func(i, j int) bool { return result[i].At < result[j].At })
	return result
}

func (a *App) AuditTrail(batchID string, filter AuditFilter) ([]domain.AuditEvent, error) {
	events, err := a.store.ListAudit(batchID)
	if err != nil {
		return nil, err
	}
	return FilterAudit(events, filter), nil
}

func (a *App) CountActions(batchID string) (map[string]int, error) {
	events, err := a.store.ListAudit(batchID)
	if err != nil {
		return nil, err
	}
	result := map[string]int{}
	for _, event := range events {
		result[event.Action]++
	}
	return result, nil
}

func (a *App) HasAuditAction(batchID, action string) (bool, error) {
	counts, err := a.CountActions(batchID)
	if err != nil {
		return false, err
	}
	return counts[action] > 0, nil
}

func (a *App) RecordHistory(recordID string) ([]domain.AuditEvent, error) {
	record, err := a.store.GetVisitor(recordID)
	if err != nil {
		return nil, err
	}
	events, err := a.store.ListAudit(record.BatchID)
	if err != nil {
		return nil, err
	}
	return FilterAudit(events, AuditFilter{RecordID: recordID}), nil
}
