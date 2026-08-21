package storage

import (
	"errors"
	"fmt"
	"path/filepath"
	"sync"
	"time"

	bolt "go.etcd.io/bbolt"
)

type Store struct {
	db   *bolt.DB
	path string
	mu   sync.RWMutex
}

func Open(path string) (*Store, error) {
	if path == "" {
		return nil, errors.New("storage path is required")
	}
	db, err := bolt.Open(filepath.Clean(path), 0600, &bolt.Options{Timeout: time.Second})
	if err != nil {
		return nil, fmt.Errorf("open storage: %w", err)
	}
	store := &Store{db: db, path: path}
	if err := store.initialize(); err != nil {
		_ = db.Close()
		return nil, err
	}
	return store, nil
}

func (s *Store) initialize() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.db.Update(func(tx *bolt.Tx) error {
		for _, name := range bucketNames {
			if _, err := tx.CreateBucketIfNotExists(name); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *Store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.db == nil {
		return nil
	}
	err := s.db.Close()
	s.db = nil
	return err
}

func (s *Store) Path() string { return s.path }

func (s *Store) view(fn func(*bolt.Tx) error) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return errors.New("storage is closed")
	}
	return s.db.View(fn)
}

func (s *Store) update(fn func(*bolt.Tx) error) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return errors.New("storage is closed")
	}
	return s.db.Update(fn)
}

func key(value string) []byte { return []byte(value) }

func requireKey(value string) error {
	if value == "" {
		return errors.New("key is required")
	}
	return nil
}

func txPut(tx *bolt.Tx, bucket, id string, value any) error {
	if err := requireKey(id); err != nil {
		return err
	}
	data, err := encode(value)
	if err != nil {
		return err
	}
	return tx.Bucket([]byte(bucket)).Put(key(id), data)
}

func txGet(tx *bolt.Tx, bucket, id string) ([]byte, error) {
	if err := requireKey(id); err != nil {
		return nil, err
	}
	raw := tx.Bucket([]byte(bucket)).Get(key(id))
	if raw == nil {
		return nil, fmt.Errorf("%s %s not found", bucket, id)
	}
	return cloneBytes(raw), nil
}

func txDelete(tx *bolt.Tx, bucket, id string) error {
	if err := requireKey(id); err != nil {
		return err
	}
	return tx.Bucket([]byte(bucket)).Delete(key(id))
}
