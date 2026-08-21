package service

import (
	"fmt"

	"parkvisitor/internal/domain"
)

func (a *App) ReopenBatch(batchID, actor string) (domain.ImportBatch, error) {
	batch, err := a.store.GetBatch(batchID)
	if err != nil {
		return batch, err
	}
	if batch.State != domain.BatchArchived {
		return batch, fmt.Errorf("only archived batch can reopen")
	}
	batch.State = domain.BatchConfirmed
	batch.ConfirmedAt = a.clock.NowString()
	if err := a.store.SaveBatch(batch); err != nil {
		return batch, err
	}
	_ = a.audit(batchID, "", "batch_reopened", actor, "archive reopened")
	return batch, nil
}

func (a *App) CountRecords(batchID string) (int, error) {
	records, err := a.store.ListVisitors(batchID)
	return len(records), err
}

func (a *App) CountPendingTasks(batchID string) (int, error) {
	tasks, err := a.store.ListTasks(batchID)
	if err != nil {
		return 0, err
	}
	count := 0
	for _, task := range tasks {
		if task.State == domain.TaskOpen {
			count++
		}
	}
	return count, nil
}
