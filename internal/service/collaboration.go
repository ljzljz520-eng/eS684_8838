package service

import (
	"errors"
	"fmt"
	"strings"

	"parkvisitor/internal/domain"
)

func (a *App) AddCollaboration(batchID, recordID, assignee, note, dueDate string) (domain.CollaborationTask, error) {
	if strings.TrimSpace(assignee) == "" {
		return domain.CollaborationTask{}, errors.New("assignee is required")
	}
	if _, err := a.store.GetBatch(batchID); err != nil {
		return domain.CollaborationTask{}, err
	}
	if recordID != "" {
		record, err := a.store.GetVisitor(recordID)
		if err != nil {
			return domain.CollaborationTask{}, err
		}
		if record.BatchID != batchID {
			return domain.CollaborationTask{}, errors.New("record does not belong to batch")
		}
	}
	task := domain.CollaborationTask{ID: a.nextID("task"), BatchID: batchID, RecordID: recordID, Assignee: assignee, State: domain.TaskOpen, Note: strings.TrimSpace(note), DueDate: dueDate}
	if err := a.store.SaveTask(task); err != nil {
		return task, err
	}
	_ = a.audit(batchID, recordID, "task_created", assignee, task.Note)
	return task, nil
}

func (a *App) CompleteTask(taskID, actor string) (domain.CollaborationTask, error) {
	task, err := a.store.GetTask(taskID)
	if err != nil {
		return task, err
	}
	if !domain.TaskStateValid(task.State) {
		return task, fmt.Errorf("invalid task state")
	}
	if task.State == domain.TaskDone {
		return task, errors.New("task already completed")
	}
	task.State = domain.TaskDone
	if err := a.store.SaveTask(task); err != nil {
		return task, err
	}
	_ = a.audit(task.BatchID, task.RecordID, "task_completed", actor, task.ID)
	return task, nil
}

func (a *App) ListCollaboration(batchID string) ([]domain.CollaborationTask, error) {
	return a.store.ListTasks(batchID)
}

func (a *App) AttachDocument(recordID, name, kind, content string) (domain.Attachment, error) {
	if strings.TrimSpace(name) == "" {
		return domain.Attachment{}, errors.New("attachment name is required")
	}
	if _, err := a.store.GetVisitor(recordID); err != nil {
		return domain.Attachment{}, err
	}
	attachment := domain.Attachment{ID: a.nextID("attachment"), RecordID: recordID, Name: name, Kind: kind, Size: int64(len(content)), Digest: fmt.Sprintf("%x", len(content))}
	if err := a.store.SaveAttachment(attachment); err != nil {
		return attachment, err
	}
	return attachment, nil
}
