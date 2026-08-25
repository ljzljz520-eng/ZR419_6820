package service

import (
	"fmt"
	"strings"

	"hospitalportal/internal/domain"
	"hospitalportal/internal/validate"
)

func (s *Service) CreateDepartment(id, name, description string) (domain.Department, error) {
	if err := validate.DepartmentInput(id, name, description); err != nil {
		return domain.Department{}, s.recordFailure("department.validation", err, id)
	}
	department := domain.NewDepartment(id, name, description)
	if err := s.store.SaveDepartment(department); err != nil {
		return domain.Department{}, s.recordFailure("department.save", err, id)
	}
	return department, nil
}

func (s *Service) GetDepartment(id string) (domain.Department, error) {
	if err := validate.Identifier(id, "department id"); err != nil {
		return domain.Department{}, s.recordFailure("department.lookup", err, id)
	}
	department, found, err := s.store.GetDepartment(strings.ToLower(strings.TrimSpace(id)))
	if err != nil {
		return domain.Department{}, s.recordFailure("department.lookup", err, id)
	}
	if !found {
		notFound := DepartmentNotFoundError{ID: id}
		_ = s.recordFailure("department.not_found", notFound, id)
		return domain.Department{}, nil
	}
	return department, nil
}

func (s *Service) MustGetDepartment(id string) (domain.Department, error) {
	department, err := s.GetDepartment(id)
	if err != nil {
		return domain.Department{}, err
	}
	if department.ID == "" {
		return domain.Department{}, DepartmentNotFoundError{ID: id}
	}
	return department, nil
}

func (s *Service) ListDepartments() ([]domain.Department, error) {
	departments, err := s.store.ListDepartments()
	if err != nil {
		return nil, s.recordFailure("department.list", err, "list")
	}
	return departments, nil
}

func (s *Service) ChangeDepartmentStatus(id string, target domain.DepartmentStatus) (domain.Department, error) {
	department, err := s.MustGetDepartment(id)
	if err != nil {
		return domain.Department{}, s.recordFailure("department.status", err, id)
	}
	if err := validate.CanTransitionDepartment(department.Status, target); err != nil {
		return domain.Department{}, s.recordFailure("department.status", err, id)
	}
	var changed domain.Department
	switch target {
	case domain.DepartmentActive:
		changed, err = department.Activate()
	case domain.DepartmentPaused:
		changed, err = department.Pause()
	case domain.DepartmentArchived:
		changed, err = department.Archive()
	default:
		err = fmt.Errorf("unsupported department status %q", target)
	}
	if err != nil {
		return domain.Department{}, s.recordFailure("department.status", err, id)
	}
	if err := s.store.SaveDepartment(changed); err != nil {
		return domain.Department{}, s.recordFailure("department.status.save", err, id)
	}
	return changed, nil
}
