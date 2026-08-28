package store

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"go.etcd.io/bbolt"
	"hospitalportal/internal/domain"
)

var bucketNames = map[string][]byte{
	"departments": []byte("departments"),
	"accounts":    []byte("accounts"),
	"dutyShifts":  []byte("dutyShifts"),
	"errorLogs":   []byte("errorLogs"),
}

type Store struct {
	db     *bbolt.DB
	path   string
	mu     sync.RWMutex
	closed bool
}

func Open(path string) (*Store, error) {
	if path == "" {
		return nil, fmt.Errorf("database path is required")
	}
	db, err := bbolt.Open(path, 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	s := &Store{db: db, path: path}
	if err := s.initialize(); err != nil {
		db.Close()
		return nil, err
	}
	return s, nil
}

func (s *Store) initialize() error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		for _, name := range bucketNames {
			if _, err := tx.CreateBucketIfNotExists(name); err != nil {
				return fmt.Errorf("create bucket: %w", err)
			}
		}
		return nil
	})
}

func (s *Store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return nil
	}
	s.closed = true
	return s.db.Close()
}

func (s *Store) Path() string {
	return s.path
}

func (s *Store) ensureOpen() error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.closed {
		return fmt.Errorf("store is closed")
	}
	return nil
}

func putJSON[T any](bucket *bbolt.Bucket, key string, value T) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return bucket.Put([]byte(key), data)
}

func getJSON[T any](bucket *bbolt.Bucket, key string, out *T) (bool, error) {
	data := bucket.Get([]byte(key))
	if data == nil {
		return false, nil
	}
	if err := json.Unmarshal(data, out); err != nil {
		return false, err
	}
	return true, nil
}

func deleteJSON(bucket *bbolt.Bucket, key string) error {
	return bucket.Delete([]byte(key))
}

func listJSON[T any](bucket *bbolt.Bucket) ([]T, error) {
	items := make([]T, 0)
	err := bucket.ForEach(func(_, value []byte) error {
		var item T
		if err := json.Unmarshal(value, &item); err != nil {
			return err
		}
		items = append(items, item)
		return nil
	})
	return items, err
}

func Remove(path string) error {
	if path == "" {
		return fmt.Errorf("database path is required")
	}
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

var _ = domain.Department{}
