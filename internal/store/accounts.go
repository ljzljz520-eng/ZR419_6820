package store

import (
	"sort"

	"go.etcd.io/bbolt"
	"hospitalportal/internal/domain"
)

func (s *Store) SaveAccount(account domain.Account) error {
	if err := s.ensureOpen(); err != nil {
		return err
	}
	if err := account.ValidateShape(); err != nil {
		return err
	}
	return s.db.Update(func(tx *bbolt.Tx) error {
		return putJSON(tx.Bucket(bucketNames["accounts"]), account.ID, account)
	})
}

func (s *Store) GetAccount(id string) (domain.Account, bool, error) {
	if err := s.ensureOpen(); err != nil {
		return domain.Account{}, false, err
	}
	var account domain.Account
	err := s.db.View(func(tx *bbolt.Tx) error {
		_, err := getJSON(tx.Bucket(bucketNames["accounts"]), id, &account)
		return err
	})
	return account, account.ID != "", err
}

func (s *Store) ListAccounts(departmentID string) ([]domain.Account, error) {
	if err := s.ensureOpen(); err != nil {
		return nil, err
	}
	var accounts []domain.Account
	err := s.db.View(func(tx *bbolt.Tx) error {
		var err error
		accounts, err = listJSON[domain.Account](tx.Bucket(bucketNames["accounts"]))
		return err
	})
	if err != nil {
		return nil, err
	}
	filtered := make([]domain.Account, 0, len(accounts))
	for _, account := range accounts {
		if departmentID == "" || account.DepartmentID == departmentID {
			filtered = append(filtered, account)
		}
	}
	sort.Slice(filtered, func(i, j int) bool { return filtered[i].ID < filtered[j].ID })
	return filtered, nil
}

func (s *Store) DeleteAccount(id string) error {
	if err := s.ensureOpen(); err != nil {
		return err
	}
	return s.db.Update(func(tx *bbolt.Tx) error {
		return deleteJSON(tx.Bucket(bucketNames["accounts"]), id)
	})
}
