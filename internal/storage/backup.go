package storage

import (
	"encoding/json"
	"fmt"
	"os"

	"parkvisitor/internal/domain"
)

type Snapshot struct {
	Batches     []domain.ImportBatch       `json:"batches"`
	Visitors    []domain.VisitorRecord     `json:"visitors"`
	Audits      []domain.AuditEvent        `json:"audits"`
	Attachments []domain.Attachment        `json:"attachments"`
	Tasks       []domain.CollaborationTask `json:"tasks"`
}

func (s *Store) Snapshot(batchID string) (Snapshot, error) {
	batches := []domain.ImportBatch{}
	batch, err := s.GetBatch(batchID)
	if err == nil {
		batches = append(batches, batch)
	}
	visitors, visitorErr := s.ListVisitors(batchID)
	if visitorErr != nil {
		return Snapshot{}, visitorErr
	}
	audits, auditErr := s.ListAudit(batchID)
	if auditErr != nil {
		return Snapshot{}, auditErr
	}
	tasks, taskErr := s.ListTasks(batchID)
	if taskErr != nil {
		return Snapshot{}, taskErr
	}
	attachments := []domain.Attachment{}
	for _, visitor := range visitors {
		items, listErr := s.ListAttachments(visitor.ID)
		if listErr != nil {
			return Snapshot{}, listErr
		}
		attachments = append(attachments, items...)
	}
	return Snapshot{Batches: batches, Visitors: visitors, Audits: audits, Attachments: attachments, Tasks: tasks}, err
}

func (s *Store) WriteSnapshot(path, batchID string) error {
	snapshot, err := s.Snapshot(batchID)
	if err != nil && len(snapshot.Batches) == 0 {
		return err
	}
	data, err := json.MarshalIndent(snapshot, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("write snapshot: %w", err)
	}
	return nil
}

func ReadSnapshot(path string) (Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Snapshot{}, err
	}
	var snapshot Snapshot
	if err := json.Unmarshal(data, &snapshot); err != nil {
		return Snapshot{}, err
	}
	return snapshot, nil
}
