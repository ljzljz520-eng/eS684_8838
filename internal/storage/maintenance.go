package storage

import (
	"errors"
	"fmt"

	bolt "go.etcd.io/bbolt"
)

func (s *Store) Exists(bucket, id string) (bool, error) {
	found := false
	err := s.view(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return fmt.Errorf("unknown bucket %s", bucket)
		}
		found = b.Get([]byte(id)) != nil
		return nil
	})
	return found, err
}

func (s *Store) ClearBucket(bucket string) error {
	return s.update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return errors.New("bucket does not exist")
		}
		cursor := b.Cursor()
		keys := [][]byte{}
		for key, _ := cursor.First(); key != nil; key, _ = cursor.Next() {
			keys = append(keys, cloneBytes(key))
		}
		for _, key := range keys {
			if err := b.Delete(key); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *Store) BucketNames() []string {
	result := make([]string, 0, len(bucketNames))
	for _, name := range bucketNames {
		result = append(result, string(name))
	}
	return result
}

func (s *Store) Health() error {
	return s.view(func(tx *bolt.Tx) error {
		for _, name := range bucketNames {
			if tx.Bucket(name) == nil {
				return fmt.Errorf("missing bucket %s", name)
			}
		}
		return nil
	})
}

func (s *Store) SaveManyVisitors(records []domainRecord) error {
	return s.update(func(tx *bolt.Tx) error {
		for _, item := range records {
			if err := txPut(tx, visitorBucket, item.ID, item.Value); err != nil {
				return err
			}
		}
		return nil
	})
}

type domainRecord struct {
	ID    string
	Value any
}
