package domain

import (
	"fmt"
	"strings"
)

type Account struct {
	ID           string
	DepartmentID string
	Username     string
	Role         AccountRole
	Active       bool
	Version      int
}

func NewAccount(id, departmentID, username string, role AccountRole) Account {
	return Account{ID: strings.TrimSpace(id), DepartmentID: strings.TrimSpace(departmentID), Username: strings.TrimSpace(username), Role: role, Active: true, Version: 1}
}

func (a Account) ValidateShape() error {
	if a.ID == "" || a.DepartmentID == "" || a.Username == "" {
		return fmt.Errorf("account %s has incomplete identity", a.ID)
	}
	if a.Version < 1 {
		return fmt.Errorf("account %s version must be positive", a.ID)
	}
	return nil
}

func (a Account) Disable() (Account, error) {
	if !a.Active {
		return a, fmt.Errorf("account %s is already inactive", a.ID)
	}
	a.Active = false
	a.Version++
	return a, nil
}

func (a Account) Enable() (Account, error) {
	if a.Active {
		return a, fmt.Errorf("account %s is already active", a.ID)
	}
	a.Active = true
	a.Version++
	return a, nil
}

func (a Account) PermissionSummary() string {
	state := "inactive"
	if a.Active {
		state = "active"
	}
	return fmt.Sprintf("%s/%s/%s", a.Username, a.Role.Label(), state)
}
