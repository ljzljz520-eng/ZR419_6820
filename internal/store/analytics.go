package store

import (
	"fmt"
	"strings"

	"hospitalportal/internal/domain"
)

type Inventory struct {
	Departments int
	Accounts    int
	Duties      int
	Errors      int
}

func (s *Store) Inventory() (Inventory, error) {
	departments, err := s.ListDepartments()
	if err != nil {
		return Inventory{}, err
	}
	accounts, err := s.ListAccounts("")
	if err != nil {
		return Inventory{}, err
	}
	duties, err := s.ListDutyShifts("")
	if err != nil {
		return Inventory{}, err
	}
	logs, err := s.ListErrorLogs()
	if err != nil {
		return Inventory{}, err
	}
	return Inventory{Departments: len(departments), Accounts: len(accounts), Duties: len(duties), Errors: len(logs)}, nil
}

func (s *Store) DepartmentSnapshot(id string) (domain.DepartmentSnapshot, error) {
	department, found, err := s.GetDepartment(id)
	if err != nil {
		return domain.DepartmentSnapshot{}, err
	}
	if !found {
		return domain.DepartmentSnapshot{}, fmt.Errorf("department %s not found", id)
	}
	accounts, err := s.ListAccounts(id)
	if err != nil {
		return domain.DepartmentSnapshot{}, err
	}
	duties, err := s.ListDutyShifts(id)
	if err != nil {
		return domain.DepartmentSnapshot{}, err
	}
	logs, err := s.ListErrorLogs()
	if err != nil {
		return domain.DepartmentSnapshot{}, err
	}
	return department.Snapshot(accounts, duties, logs), nil
}

func (s *Store) FindErrorLogsByContext(context string) ([]domain.ErrorLog, error) {
	logs, err := s.ListErrorLogs()
	if err != nil {
		return nil, err
	}
	context = strings.TrimSpace(context)
	filtered := make([]domain.ErrorLog, 0, len(logs))
	for _, log := range logs {
		if context == "" || strings.Contains(log.Context, context) {
			filtered = append(filtered, log)
		}
	}
	return filtered, nil
}

func (s *Store) ReplaceDepartment(department domain.Department) (domain.Department, error) {
	if err := s.SaveDepartment(department); err != nil {
		return domain.Department{}, err
	}
	stored, found, err := s.GetDepartment(department.ID)
	if err != nil {
		return domain.Department{}, err
	}
	if !found {
		return domain.Department{}, fmt.Errorf("department %s disappeared after save", department.ID)
	}
	return stored, nil
}
