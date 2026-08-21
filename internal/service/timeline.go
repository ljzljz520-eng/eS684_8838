package service

import "sort"

type TimelineItem struct {
	At       string `json:"at"`
	Action   string `json:"action"`
	Actor    string `json:"actor"`
	Detail   string `json:"detail"`
	RecordID string `json:"record_id"`
}

func (a *App) Timeline(batchID string) ([]TimelineItem, error) {
	events, err := a.store.ListAudit(batchID)
	if err != nil {
		return nil, err
	}
	result := make([]TimelineItem, 0, len(events))
	for _, event := range events {
		result = append(result, TimelineItem{At: event.At, Action: event.Action, Actor: event.Actor, Detail: event.Detail, RecordID: event.RecordID})
	}
	sort.SliceStable(result, func(i, j int) bool { return result[i].At < result[j].At })
	return result, nil
}

func (a *App) LastAction(batchID string) (TimelineItem, error) {
	timeline, err := a.Timeline(batchID)
	if err != nil {
		return TimelineItem{}, err
	}
	if len(timeline) == 0 {
		return TimelineItem{}, domainErrorForBatch(batchID)
	}
	return timeline[len(timeline)-1], nil
}

type batchError struct{ id string }

func (e *batchError) Error() string       { return "no audit events for batch " + e.id }
func domainErrorForBatch(id string) error { return &batchError{id: id} }
