package validate

import (
	"fmt"
	"regexp"
	"strings"

	"hospitalportal/internal/domain"
)

var identifierPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_-]{1,31}$`)

func Identifier(value, label string) error {
	value = strings.TrimSpace(value)
	if value == "" {
		return fmt.Errorf("%s is required", label)
	}
	if !identifierPattern.MatchString(value) {
		return fmt.Errorf("%s %q has invalid format", label, value)
	}
	return nil
}

func DepartmentInput(id, name, description string) error {
	if err := Identifier(id, "department id"); err != nil {
		return err
	}
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("department name is required")
	}
	if len([]rune(description)) > 240 {
		return fmt.Errorf("department description is too long")
	}
	return nil
}

func AccountInput(id, departmentID, username string, role domain.AccountRole) error {
	if err := Identifier(id, "account id"); err != nil {
		return err
	}
	if err := Identifier(departmentID, "department id"); err != nil {
		return err
	}
	if strings.TrimSpace(username) == "" {
		return fmt.Errorf("username is required")
	}
	if !IsRoleSupported(role) {
		return fmt.Errorf("unsupported role %q", role)
	}
	return nil
}

func DutyInput(id, departmentID, dateKey, clinician string, shift domain.ShiftName) error {
	if err := Identifier(id, "duty shift id"); err != nil {
		return err
	}
	if err := Identifier(departmentID, "department id"); err != nil {
		return err
	}
	if !validDateKey(dateKey) {
		return fmt.Errorf("date key %q must use YYYY-MM-DD", dateKey)
	}
	if strings.TrimSpace(clinician) == "" {
		return fmt.Errorf("clinician is required")
	}
	if !IsShiftSupported(shift) {
		return fmt.Errorf("unsupported shift %q", shift)
	}
	return nil
}

func IsRoleSupported(role domain.AccountRole) bool {
	switch role {
	case domain.RoleDoctor, domain.RoleNurse, domain.RoleAdministrator:
		return true
	default:
		return false
	}
}

func IsShiftSupported(shift domain.ShiftName) bool {
	switch shift {
	case domain.ShiftMorning, domain.ShiftEvening, domain.ShiftNight:
		return true
	default:
		return false
	}
}

func validDateKey(value string) bool {
	if len(value) != 10 || value[4] != '-' || value[7] != '-' {
		return false
	}
	for i, r := range value {
		if i == 4 || i == 7 {
			continue
		}
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}
