package store

import (
	"fmt"
	"sort"

	"go.etcd.io/bbolt"
	"hospitalportal/internal/domain"
)

func (s *Store) SaveDepartment(department domain.Department) error {
	if err := s.ensureOpen(); err != nil {
		return err
	}
	if err := department.ValidateShape(); err != nil {
		return err
	}
	return s.db.Update(func(tx *bbolt.Tx) error {
		return putJSON(tx.Bucket(bucketNames["departments"]), department.ID, department)
	})
}

func (s *Store) GetDepartment(id string) (domain.Department, bool, error) {
	if err := s.ensureOpen(); err != nil {
		return domain.Department{}, false, err
	}
	var department domain.Department
	err := s.db.View(func(tx *bbolt.Tx) error {
		found, err := getJSON(tx.Bucket(bucketNames["departments"]), id, &department)
		if err != nil {
			return err
		}
		if !found {
			return nil
		}
		return nil
	})
	return department, department.ID != "", err
}

func (s *Store) ListDepartments() ([]domain.Department, error) {
	if err := s.ensureOpen(); err != nil {
		return nil, err
	}
	var departments []domain.Department
	err := s.db.View(func(tx *bbolt.Tx) error {
		var err error
		departments, err = listJSON[domain.Department](tx.Bucket(bucketNames["departments"]))
		return err
	})
	if err != nil {
		return nil, err
	}
	sort.Slice(departments, func(i, j int) bool { return departments[i].ID < departments[j].ID })
	return departments, nil
}

func (s *Store) DeleteDepartment(id string) error {
	if err := s.ensureOpen(); err != nil {
		return err
	}
	if id == "" {
		return fmt.Errorf("department id is required")
	}
	return s.db.Update(func(tx *bbolt.Tx) error {
		return deleteJSON(tx.Bucket(bucketNames["departments"]), id)
	})
}
