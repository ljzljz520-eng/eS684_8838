package service

import (
	"errors"
	"sort"
	"strings"

	"parkvisitor/internal/domain"
)

type ReviewQueue struct {
	Records  []domain.VisitorRecord `json:"records"`
	Total    int                    `json:"total"`
	Priority string                 `json:"priority"`
}

func (a *App) BuildReviewQueue(batchID string) (ReviewQueue, error) {
	records, err := a.store.ListVisitors(batchID)
	if err != nil {
		return ReviewQueue{}, err
	}
	queue := ReviewQueue{Records: []domain.VisitorRecord{}, Priority: "normal"}
	for _, record := range records {
		if record.Status == domain.StatusNeedsReview {
			queue.Records = append(queue.Records, record)
		}
	}
	sort.Slice(queue.Records, func(i, j int) bool { return queue.Records[i].VisitDate < queue.Records[j].VisitDate })
	queue.Total = len(queue.Records)
	if queue.Total > 10 {
		queue.Priority = "high"
	}
	return queue, nil
}

func (a *App) AssignReviewQueue(batchID, actor string) (ReviewQueue, error) {
	if strings.TrimSpace(actor) == "" {
		return ReviewQueue{}, errors.New("reviewer is required")
	}
	queue, err := a.BuildReviewQueue(batchID)
	if err != nil {
		return queue, err
	}
	for _, record := range queue.Records {
		if _, taskErr := a.AddCollaboration(batchID, record.ID, actor, "review visitor", record.VisitDate); taskErr != nil {
			return ReviewQueue{}, taskErr
		}
	}
	return queue, nil
}

func (a *App) QueueSize(batchID string) (int, error) {
	queue, err := a.BuildReviewQueue(batchID)
	return queue.Total, err
}

func (a *App) HasReviewWork(batchID string) (bool, error) {
	size, err := a.QueueSize(batchID)
	return size > 0, err
}
