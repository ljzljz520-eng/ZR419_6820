package service

import (
	"errors"
	"fmt"

	"hospitalportal/internal/audit"
	"hospitalportal/internal/domain"
	"hospitalportal/internal/store"
	"hospitalportal/internal/validate"
)

var ErrDepartmentNotFound = errors.New("department not found")
var ErrAccountNotFound = errors.New("account not found")
var ErrDutyNotFound = errors.New("duty shift not found")

type DepartmentNotFoundError struct {
	ID string
}

func (e DepartmentNotFoundError) Error() string {
	return fmt.Sprintf("department %s not found", e.ID)
}

func (e DepartmentNotFoundError) Unwrap() error {
	return ErrDepartmentNotFound
}

type AccountNotFoundError struct {
	ID string
}

func (e AccountNotFoundError) Error() string {
	return fmt.Sprintf("account %s not found", e.ID)
}

func (e AccountNotFoundError) Unwrap() error {
	return ErrAccountNotFound
}

type DutyNotFoundError struct {
	ID string
}

func (e DutyNotFoundError) Error() string {
	return fmt.Sprintf("duty shift %s not found", e.ID)
}

func (e DutyNotFoundError) Unwrap() error {
	return ErrDutyNotFound
}

type Service struct {
	store *store.Store
	audit *audit.Logger
}

func New(s *store.Store) *Service {
	return &Service{store: s, audit: audit.NewLogger(s)}
}

func (s *Service) Store() *store.Store {
	return s.store
}

func (s *Service) Audit() *audit.Logger {
	return s.audit
}

func (s *Service) recordFailure(code string, err error, context string) error {
	if err == nil {
		return nil
	}
	_, logErr := s.audit.Record(code, err.Error(), context)
	if logErr != nil {
		return fmt.Errorf("%w; record audit: %v", err, logErr)
	}
	return err
}

func normalizeDepartment(d domain.Department) domain.Department {
	d.ID = validate.NormalizeUsername(d.ID)
	d.Name = validate.NormalizeUsername(d.Name)
	return d
}
