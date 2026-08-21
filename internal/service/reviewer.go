package service

import (
	"errors"
	"fmt"

	"parkvisitor/internal/domain"
)

func (a *App) ReviewBatch(batchID, actor string) ([]domain.VisitorRecord, error) {
	batch, err := a.store.GetBatch(batchID)
	if err != nil {
		return nil, err
	}
	if !domain.CanReview(batch.State) {
		return nil, fmt.Errorf("batch %s cannot be reviewed in state %s", batchID, batch.State)
	}
	records, err := a.store.ListVisitors(batchID)
	if err != nil {
		return nil, err
	}
	for i := range records {
		if records[i].Status == domain.StatusValidated {
			records[i].Status = domain.StatusApproved
			records[i].UpdatedAt = a.clock.NowString()
			if err := a.store.SaveVisitor(records[i]); err != nil {
				return nil, err
			}
			_ = a.audit(batchID, records[i].ID, "approved", actor, "review passed")
		}
	}
	batch.State = domain.BatchConfirmed
	batch.ConfirmedAt = a.clock.NowString()
	if err := a.store.SaveBatch(batch); err != nil {
		return nil, err
	}
	return records, nil
}

func (a *App) ApproveRecord(batchID, recordID, actor string) (domain.VisitorRecord, error) {
	record, err := a.store.GetVisitor(recordID)
	if err != nil {
		return record, err
	}
	if record.BatchID != batchID {
		return record, errors.New("record does not belong to batch")
	}
	if record.Status != domain.StatusValidated && record.Status != domain.StatusNeedsReview {
		return record, fmt.Errorf("record cannot be approved from %s", record.Status)
	}
	record.Status = domain.StatusApproved
	record.UpdatedAt = a.clock.NowString()
	if err := a.store.SaveVisitor(record); err != nil {
		return record, err
	}
	_ = a.audit(batchID, recordID, "approved", actor, "manual approval")
	return record, nil
}

func (a *App) RejectRecord(batchID, recordID, actor, reason string) (domain.VisitorRecord, error) {
	record, err := a.store.GetVisitor(recordID)
	if err != nil {
		return record, err
	}
	if record.BatchID != batchID {
		return record, errors.New("record does not belong to batch")
	}
	if reason == "" {
		return record, errors.New("rejection reason is required")
	}
	record.Status = domain.StatusRejected
	record.Notes = reason
	record.UpdatedAt = a.clock.NowString()
	if err := a.store.SaveVisitor(record); err != nil {
		return record, err
	}
	_ = a.audit(batchID, recordID, "rejected", actor, reason)
	return record, nil
}

func (a *App) PublishBatch(batchID, actor string) (domain.ImportBatch, error) {
	batch, err := a.store.GetBatch(batchID)
	if err != nil {
		return batch, err
	}
	if !domain.CanPublish(batch.State) {
		return batch, fmt.Errorf("batch cannot be published from %s", batch.State)
	}
	records, err := a.store.ListVisitors(batchID)
	if err != nil {
		return batch, err
	}
	for i := range records {
		if records[i].Status == domain.StatusApproved || records[i].Status == domain.StatusValidated {
			records[i].Status = domain.StatusPublished
			records[i].UpdatedAt = a.clock.NowString()
			if err := a.store.SaveVisitor(records[i]); err != nil {
				return batch, err
			}
		}
	}
	batch.State = domain.BatchPublished
	if err := a.store.SaveBatch(batch); err != nil {
		return batch, err
	}
	_ = a.audit(batchID, "", "published", actor, fmt.Sprintf("records=%d", len(records)))
	return batch, nil
}
