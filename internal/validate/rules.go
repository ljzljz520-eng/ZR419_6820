package validate

import (
	"fmt"
	"strings"

	"hospitalportal/internal/domain"
)

func CanCreateAccount(department domain.Department, role domain.AccountRole) error {
	if !department.Status.CanAcceptAccounts() {
		return fmt.Errorf("department %s is not accepting accounts", department.ID)
	}
	if !IsRoleSupported(role) {
		return fmt.Errorf("unsupported role %q", role)
	}
	return nil
}

func CanPublishDuty(department domain.Department, clinician string) error {
	if !department.Status.IsVisible() {
		return fmt.Errorf("department %s is not visible", department.ID)
	}
	if strings.TrimSpace(clinician) == "" {
		return fmt.Errorf("clinician is required")
	}
	return nil
}

func HasDutyConflict(existing []domain.DutyShift, candidate domain.DutyShift) bool {
	for _, duty := range existing {
		if duty.ID == candidate.ID {
			continue
		}
		if duty.Status == domain.RecordCancelled {
			continue
		}
		if duty.SlotKey() == candidate.SlotKey() {
			return true
		}
	}
	return false
}

func CanTransitionDepartment(current, next domain.DepartmentStatus) error {
	if current == domain.DepartmentArchived && next != domain.DepartmentArchived {
		return fmt.Errorf("archived department cannot transition")
	}
	if current == next {
		return fmt.Errorf("department is already %s", current)
	}
	return nil
}

func NormalizeUsername(username string) string {
	return strings.ToLower(strings.TrimSpace(username))
}
