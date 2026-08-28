package store

import (
	"sort"

	"go.etcd.io/bbolt"
	"hospitalportal/internal/domain"
)

func (s *Store) SaveErrorLog(log domain.ErrorLog) error {
	if err := s.ensureOpen(); err != nil {
		return err
	}
	if err := log.ValidateShape(); err != nil {
		return err
	}
	return s.db.Update(func(tx *bbolt.Tx) error {
		return putJSON(tx.Bucket(bucketNames["errorLogs"]), log.ID, log)
	})
}

func (s *Store) ListErrorLogs() ([]domain.ErrorLog, error) {
	if err := s.ensureOpen(); err != nil {
		return nil, err
	}
	var logs []domain.ErrorLog
	err := s.db.View(func(tx *bbolt.Tx) error {
		var err error
		logs, err = listJSON[domain.ErrorLog](tx.Bucket(bucketNames["errorLogs"]))
		return err
	})
	if err != nil {
		return nil, err
	}
	sort.Slice(logs, func(i, j int) bool { return logs[i].ID < logs[j].ID })
	return logs, nil
}

func (s *Store) CountErrors() (int, error) {
	logs, err := s.ListErrorLogs()
	return len(logs), err
}
