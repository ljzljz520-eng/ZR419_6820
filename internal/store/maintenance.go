package store

import (
	"fmt"
	"sort"
	"strings"

	"hospitalportal/internal/domain"
)

type MaintenanceResult struct {
	UpdatedDepartments int
	UpdatedAccounts    int
	UpdatedDuties      int
	Messages           []string
}

func (s *Store) NormalizeRecords() (MaintenanceResult, error) {
	departments, err := s.ListDepartments()
	if err != nil {
		return MaintenanceResult{}, err
	}
	accounts, err := s.ListAccounts("")
	if err != nil {
		return MaintenanceResult{}, err
	}
	duties, err := s.ListDutyShifts("")
	if err != nil {
		return MaintenanceResult{}, err
	}
	result := MaintenanceResult{Messages: make([]string, 0)}
	for _, department := range departments {
		normalized := department
		normalized.ID = strings.ToLower(strings.TrimSpace(normalized.ID))
		normalized.Name = strings.TrimSpace(normalized.Name)
		if normalized.ID != department.ID || normalized.Name != department.Name {
			if err := s.SaveDepartment(normalized); err != nil {
				return MaintenanceResult{}, err
			}
			result.UpdatedDepartments++
		}
	}
	for _, account := range accounts {
		normalized := account
		normalized.Username = strings.ToLower(strings.TrimSpace(normalized.Username))
		if normalized.Username != account.Username {
			if err := s.SaveAccount(normalized); err != nil {
				return MaintenanceResult{}, err
			}
			result.UpdatedAccounts++
		}
	}
	for _, duty := range duties {
		normalized := duty
		normalized.DateKey = strings.TrimSpace(normalized.DateKey)
		normalized.Clinician = strings.TrimSpace(normalized.Clinician)
		if normalized.DateKey != duty.DateKey || normalized.Clinician != duty.Clinician {
			if err := s.SaveDutyShift(normalized); err != nil {
				return MaintenanceResult{}, err
			}
			result.UpdatedDuties++
		}
	}
	result.Messages = append(result.Messages, fmt.Sprintf("normalized %d departments", result.UpdatedDepartments))
	result.Messages = append(result.Messages, fmt.Sprintf("normalized %d accounts", result.UpdatedAccounts))
	result.Messages = append(result.Messages, fmt.Sprintf("normalized %d duty shifts", result.UpdatedDuties))
	return result, nil
}

func (s *Store) DepartmentIDs() ([]string, error) {
	departments, err := s.ListDepartments()
	if err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(departments))
	for _, department := range departments {
		ids = append(ids, department.ID)
	}
	sort.Strings(ids)
	return ids, nil
}

func (s *Store) ValidateAll() []error {
	errors := make([]error, 0)
	departments, err := s.ListDepartments()
	if err != nil {
		return append(errors, err)
	}
	for _, department := range departments {
		if err := department.ValidateShape(); err != nil {
			errors = append(errors, err)
		}
	}
	accounts, err := s.ListAccounts("")
	if err != nil {
		return append(errors, err)
	}
	for _, account := range accounts {
		if err := account.ValidateShape(); err != nil {
			errors = append(errors, err)
		}
	}
	duties, err := s.ListDutyShifts("")
	if err != nil {
		return append(errors, err)
	}
	for _, duty := range duties {
		if err := duty.ValidateShape(); err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}

func (s *Store) SeedReferenceData() error {
	references := []domain.Department{
		domain.NewDepartment("emergency", "Emergency", "Urgent care"),
		domain.NewDepartment("pediatrics", "Pediatrics", "Children care"),
	}
	for _, department := range references {
		if _, found, err := s.GetDepartment(department.ID); err != nil {
			return err
		} else if !found {
			if err := s.SaveDepartment(department); err != nil {
				return err
			}
		}
	}
	return nil
}
