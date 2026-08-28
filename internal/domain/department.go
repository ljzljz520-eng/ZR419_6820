package domain

import (
	"fmt"
	"strings"
)

type Department struct {
	ID          string
	Name        string
	Description string
	Status      DepartmentStatus
	Version     int
}

func NewDepartment(id, name, description string) Department {
	return Department{ID: strings.TrimSpace(id), Name: strings.TrimSpace(name), Description: strings.TrimSpace(description), Status: DepartmentActive, Version: 1}
}

func (d Department) ValidateShape() error {
	if strings.TrimSpace(d.ID) == "" {
		return fmt.Errorf("department id is required")
	}
	if strings.TrimSpace(d.Name) == "" {
		return fmt.Errorf("department %s name is required", d.ID)
	}
	if d.Version < 1 {
		return fmt.Errorf("department %s version must be positive", d.ID)
	}
	return nil
}

func (d Department) Pause() (Department, error) {
	if d.Status == DepartmentArchived {
		return d, fmt.Errorf("department %s is archived", d.ID)
	}
	d.Status = DepartmentPaused
	d.Version++
	return d, nil
}

func (d Department) Activate() (Department, error) {
	if d.Status == DepartmentArchived {
		return d, fmt.Errorf("department %s is archived", d.ID)
	}
	d.Status = DepartmentActive
	d.Version++
	return d, nil
}

func (d Department) Archive() (Department, error) {
	if d.Status == DepartmentArchived {
		return d, fmt.Errorf("department %s is already archived", d.ID)
	}
	d.Status = DepartmentArchived
	d.Version++
	return d, nil
}

func (d Department) Summary() string {
	return fmt.Sprintf("%s | %s | %s", d.ID, d.Name, d.Status)
}
