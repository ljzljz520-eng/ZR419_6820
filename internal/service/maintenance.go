package service

import (
	"fmt"

	"hospitalportal/internal/domain"
	"hospitalportal/internal/report"
)

func (s *Service) MaintenanceSummary() (string, error) {
	inventory, err := s.store.Inventory()
	if err != nil {
		return "", s.recordFailure("maintenance.inventory", err, "global")
	}
	return fmt.Sprintf("departments=%d accounts=%d duties=%d errors=%d", inventory.Departments, inventory.Accounts, inventory.Duties, inventory.Errors), nil
}

func (s *Service) PermissionMatrix(departmentID string) (report.PermissionMatrix, error) {
	department, err := s.MustGetDepartment(departmentID)
	if err != nil {
		return report.PermissionMatrix{}, err
	}
	accounts, err := s.store.ListAccounts(departmentID)
	if err != nil {
		return report.PermissionMatrix{}, s.recordFailure("permissions.accounts", err, departmentID)
	}
	return report.BuildPermissionMatrix(department, accounts), nil
}

func (s *Service) NormalizeRecords() (string, error) {
	result, err := s.store.NormalizeRecords()
	if err != nil {
		return "", s.recordFailure("maintenance.normalize", err, "global")
	}
	return fmt.Sprintf("departments=%d accounts=%d duties=%d", result.UpdatedDepartments, result.UpdatedAccounts, result.UpdatedDuties), nil
}

func (s *Service) ValidateData() []error {
	return s.store.ValidateAll()
}

func (s *Service) SeedDefaults() error {
	if err := s.store.SeedReferenceData(); err != nil {
		return s.recordFailure("maintenance.seed", err, "global")
	}
	return nil
}

func (s *Service) DepartmentSnapshots() ([]domain.DepartmentSnapshot, error) {
	departments, err := s.ListDepartments()
	if err != nil {
		return nil, err
	}
	snapshots := make([]domain.DepartmentSnapshot, 0, len(departments))
	for _, department := range departments {
		snapshot, err := s.store.DepartmentSnapshot(department.ID)
		if err != nil {
			return nil, err
		}
		snapshots = append(snapshots, snapshot)
	}
	return snapshots, nil
}
