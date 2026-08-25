package service

import (
	"fmt"

	"hospitalportal/internal/domain"
	"hospitalportal/internal/validate"
)

func (s *Service) CreateAccount(id, departmentID, username string, role domain.AccountRole) (domain.Account, error) {
	if err := validate.AccountInput(id, departmentID, username, role); err != nil {
		return domain.Account{}, s.recordFailure("account.validation", err, id)
	}
	department, err := s.MustGetDepartment(departmentID)
	if err != nil {
		return domain.Account{}, s.recordFailure("account.department", err, departmentID)
	}
	if err := validate.CanCreateAccount(department, role); err != nil {
		return domain.Account{}, s.recordFailure("account.permission", err, id)
	}
	account := domain.NewAccount(id, departmentID, validate.NormalizeUsername(username), role)
	if err := s.store.SaveAccount(account); err != nil {
		return domain.Account{}, s.recordFailure("account.save", err, id)
	}
	return account, nil
}

func (s *Service) GetAccount(id string) (domain.Account, error) {
	account, found, err := s.store.GetAccount(id)
	if err != nil {
		return domain.Account{}, s.recordFailure("account.lookup", err, id)
	}
	if !found {
		notFound := AccountNotFoundError{ID: id}
		return domain.Account{}, s.recordFailure("account.not_found", notFound, id)
	}
	return account, nil
}

func (s *Service) ListAccounts(departmentID string) ([]domain.Account, error) {
	if departmentID != "" {
		if err := validate.Identifier(departmentID, "department id"); err != nil {
			return nil, s.recordFailure("account.list", err, departmentID)
		}
	}
	accounts, err := s.store.ListAccounts(departmentID)
	if err != nil {
		return nil, s.recordFailure("account.list", err, departmentID)
	}
	return accounts, nil
}

func (s *Service) SetAccountActive(id string, active bool) (domain.Account, error) {
	account, err := s.GetAccount(id)
	if err != nil {
		return domain.Account{}, err
	}
	if account.Active == active {
		return account, fmt.Errorf("account %s already has active=%t", id, active)
	}
	var changed domain.Account
	if active {
		changed, err = account.Enable()
	} else {
		changed, err = account.Disable()
	}
	if err != nil {
		return domain.Account{}, s.recordFailure("account.status", err, id)
	}
	if err := s.store.SaveAccount(changed); err != nil {
		return domain.Account{}, s.recordFailure("account.status.save", err, id)
	}
	return changed, nil
}
