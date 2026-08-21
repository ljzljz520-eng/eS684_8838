package service

import (
	"fmt"

	"parkvisitor/internal/domain"
)

func (a *App) ArchiveBatch(batchID, actor string) (domain.ImportBatch, error) {
	batch, err := a.store.GetBatch(batchID)
	if err != nil {
		return batch, err
	}
	if !domain.CanArchive(batch.State) {
		return batch, fmt.Errorf("batch cannot be archived from %s", batch.State)
	}
	records, err := a.store.ListVisitors(batchID)
	if err != nil {
		return batch, err
	}
	for i := range records {
		records[i].Status = domain.StatusArchived
		records[i].UpdatedAt = a.clock.NowString()
		if err := a.store.SaveVisitor(records[i]); err != nil {
			return batch, err
		}
	}
	batch.State = domain.BatchArchived
	if err := a.store.SaveBatch(batch); err != nil {
		return batch, err
	}
	_ = a.audit(batchID, "", "archived", actor, fmt.Sprintf("records=%d", len(records)))
	return batch, nil
}

func (a *App) RestoreVisitor(recordID, actor string) (domain.VisitorRecord, error) {
	record, err := a.store.GetVisitor(recordID)
	if err != nil {
		return record, err
	}
	if record.Status != domain.StatusArchived {
		return record, fmt.Errorf("record is not archived")
	}
	record.Status = domain.StatusPublished
	record.UpdatedAt = a.clock.NowString()
	if err := a.store.SaveVisitor(record); err != nil {
		return record, err
	}
	_ = a.audit(record.BatchID, record.ID, "restored", actor, "restored from archive")
	return record, nil
}

func (a *App) PurgeVisitor(recordID string) error {
	record, err := a.store.GetVisitor(recordID)
	if err != nil {
		return err
	}
	if record.Status != domain.StatusArchived {
		return fmt.Errorf("only archived records can be purged")
	}
	return a.store.RemoveVisitor(record.ID)
}
