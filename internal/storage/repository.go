package storage

import (
	"fmt"
	"sort"

	bolt "go.etcd.io/bbolt"
	"parkvisitor/internal/domain"
)

func (s *Store) SaveVisitor(record domain.VisitorRecord) error {
	return s.update(func(tx *bolt.Tx) error { return txPut(tx, visitorBucket, record.ID, record) })
}

func (s *Store) GetVisitor(id string) (domain.VisitorRecord, error) {
	var result domain.VisitorRecord
	err := s.view(func(tx *bolt.Tx) error {
		raw, err := txGet(tx, visitorBucket, id)
		if err != nil {
			return err
		}
		return decode(raw, &result)
	})
	return result, err
}

func (s *Store) ListVisitors(batchID string) ([]domain.VisitorRecord, error) {
	result := []domain.VisitorRecord{}
	err := s.view(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte(visitorBucket)).ForEach(func(_, raw []byte) error {
			var item domain.VisitorRecord
			if err := decode(raw, &item); err != nil {
				return err
			}
			if batchID == "" || item.BatchID == batchID {
				result = append(result, item)
			}
			return nil
		})
	})
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result, err
}

func (s *Store) SaveBatch(batch domain.ImportBatch) error {
	return s.update(func(tx *bolt.Tx) error { return txPut(tx, batchBucket, batch.ID, batch) })
}

func (s *Store) GetBatch(id string) (domain.ImportBatch, error) {
	var result domain.ImportBatch
	err := s.view(func(tx *bolt.Tx) error {
		raw, err := txGet(tx, batchBucket, id)
		if err != nil {
			return err
		}
		return decode(raw, &result)
	})
	return result, err
}

func (s *Store) SaveAudit(event domain.AuditEvent) error {
	return s.update(func(tx *bolt.Tx) error { return txPut(tx, auditBucket, event.ID, event) })
}

func (s *Store) ListAudit(batchID string) ([]domain.AuditEvent, error) {
	result := []domain.AuditEvent{}
	err := s.view(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte(auditBucket)).ForEach(func(_, raw []byte) error {
			var item domain.AuditEvent
			if err := decode(raw, &item); err != nil {
				return err
			}
			if batchID == "" || item.BatchID == batchID {
				result = append(result, item)
			}
			return nil
		})
	})
	sort.Slice(result, func(i, j int) bool { return result[i].At < result[j].At })
	return result, err
}

func (s *Store) SaveAttachment(item domain.Attachment) error {
	return s.update(func(tx *bolt.Tx) error { return txPut(tx, attachmentBucket, item.ID, item) })
}

func (s *Store) ListAttachments(recordID string) ([]domain.Attachment, error) {
	result := []domain.Attachment{}
	err := s.view(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte(attachmentBucket)).ForEach(func(_, raw []byte) error {
			var item domain.Attachment
			if err := decode(raw, &item); err != nil {
				return err
			}
			if recordID == "" || item.RecordID == recordID {
				result = append(result, item)
			}
			return nil
		})
	})
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result, err
}

func (s *Store) SaveTask(item domain.CollaborationTask) error {
	return s.update(func(tx *bolt.Tx) error { return txPut(tx, taskBucket, item.ID, item) })
}

func (s *Store) GetTask(id string) (domain.CollaborationTask, error) {
	var result domain.CollaborationTask
	err := s.view(func(tx *bolt.Tx) error {
		raw, err := txGet(tx, taskBucket, id)
		if err != nil {
			return err
		}
		return decode(raw, &result)
	})
	return result, err
}

func (s *Store) ListTasks(batchID string) ([]domain.CollaborationTask, error) {
	result := []domain.CollaborationTask{}
	err := s.view(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte(taskBucket)).ForEach(func(_, raw []byte) error {
			var item domain.CollaborationTask
			if err := decode(raw, &item); err != nil {
				return err
			}
			if batchID == "" || item.BatchID == batchID {
				result = append(result, item)
			}
			return nil
		})
	})
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result, err
}

func (s *Store) RemoveVisitor(id string) error {
	return s.update(func(tx *bolt.Tx) error { return txDelete(tx, visitorBucket, id) })
}

func (s *Store) Count(bucket string) (int, error) {
	count := 0
	err := s.view(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return fmt.Errorf("unknown bucket %s", bucket)
		}
		count = b.Stats().KeyN
		return nil
	})
	return count, err
}
