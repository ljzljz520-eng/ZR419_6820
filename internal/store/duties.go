package store

import (
	"sort"

	"go.etcd.io/bbolt"
	"hospitalportal/internal/domain"
)

func (s *Store) SaveDutyShift(duty domain.DutyShift) error {
	if err := s.ensureOpen(); err != nil {
		return err
	}
	if err := duty.ValidateShape(); err != nil {
		return err
	}
	return s.db.Update(func(tx *bbolt.Tx) error {
		return putJSON(tx.Bucket(bucketNames["dutyShifts"]), duty.ID, duty)
	})
}

func (s *Store) GetDutyShift(id string) (domain.DutyShift, bool, error) {
	if err := s.ensureOpen(); err != nil {
		return domain.DutyShift{}, false, err
	}
	var duty domain.DutyShift
	err := s.db.View(func(tx *bbolt.Tx) error {
		_, err := getJSON(tx.Bucket(bucketNames["dutyShifts"]), id, &duty)
		return err
	})
	return duty, duty.ID != "", err
}

func (s *Store) ListDutyShifts(departmentID string) ([]domain.DutyShift, error) {
	if err := s.ensureOpen(); err != nil {
		return nil, err
	}
	var duties []domain.DutyShift
	err := s.db.View(func(tx *bbolt.Tx) error {
		var err error
		duties, err = listJSON[domain.DutyShift](tx.Bucket(bucketNames["dutyShifts"]))
		return err
	})
	if err != nil {
		return nil, err
	}
	filtered := make([]domain.DutyShift, 0, len(duties))
	for _, duty := range duties {
		if departmentID == "" || duty.DepartmentID == departmentID {
			filtered = append(filtered, duty)
		}
	}
	sort.Slice(filtered, func(i, j int) bool {
		if filtered[i].DateKey == filtered[j].DateKey {
			return filtered[i].Shift < filtered[j].Shift
		}
		return filtered[i].DateKey < filtered[j].DateKey
	})
	return filtered, nil
}

func (s *Store) DeleteDutyShift(id string) error {
	if err := s.ensureOpen(); err != nil {
		return err
	}
	return s.db.Update(func(tx *bbolt.Tx) error {
		return deleteJSON(tx.Bucket(bucketNames["dutyShifts"]), id)
	})
}
