package api

import (
	"fmt"
	"strings"

	"hospitalportal/internal/domain"
	"hospitalportal/internal/service"
)

type DepartmentRequest struct {
	ID          string
	Name        string
	Description string
}

type AccountRequest struct {
	ID           string
	DepartmentID string
	Username     string
	Role         domain.AccountRole
}

type DutyRequest struct {
	ID           string
	DepartmentID string
	DateKey      string
	Clinician    string
	Shift        domain.ShiftName
}

type Result struct {
	OK      bool
	Kind    string
	Message string
	Rows    []string
	Error   string
}

type Handler struct {
	service *service.Service
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) CreateDepartment(request DepartmentRequest) Result {
	department, err := h.service.CreateDepartment(request.ID, request.Name, request.Description)
	if err != nil {
		return failure("department", err)
	}
	return success("department", fmt.Sprintf("department %s saved", department.ID), []string{department.Summary()})
}

func (h *Handler) FindDepartment(id string) Result {
	department, err := h.service.MustGetDepartment(id)
	if err != nil {
		return failure("department", err)
	}
	return success("department", "department found", []string{department.Summary(), department.Description})
}

func (h *Handler) CreateAccount(request AccountRequest) Result {
	account, err := h.service.CreateAccount(request.ID, request.DepartmentID, request.Username, request.Role)
	if err != nil {
		return failure("account", err)
	}
	return success("account", "account saved", []string{account.PermissionSummary()})
}

func (h *Handler) CreateDuty(request DutyRequest) Result {
	duty, err := h.service.CreateDutyShift(request.ID, request.DepartmentID, request.DateKey, request.Clinician, request.Shift)
	if err != nil {
		return failure("duty", err)
	}
	return success("duty", "duty saved", []string{formatDuty(duty)})
}

func (h *Handler) PublishDuty(id string) Result {
	duty, err := h.service.PublishDutyShift(id)
	if err != nil {
		return failure("duty", err)
	}
	return success("duty", "duty published", []string{formatDuty(duty)})
}

func (h *Handler) ListDepartmentRows() Result {
	departments, err := h.service.ListDepartments()
	if err != nil {
		return failure("department", err)
	}
	rows := make([]string, 0, len(departments))
	for _, department := range departments {
		rows = append(rows, department.Summary())
	}
	return success("department", fmt.Sprintf("%d departments", len(rows)), rows)
}

func (h *Handler) ListDutyRows(departmentID string) Result {
	duties, err := h.service.ListDutyShifts(departmentID)
	if err != nil {
		return failure("duty", err)
	}
	rows := make([]string, 0, len(duties))
	for _, duty := range duties {
		rows = append(rows, formatDuty(duty))
	}
	return success("duty", fmt.Sprintf("%d duty shifts", len(rows)), rows)
}

func success(kind, message string, rows []string) Result {
	return Result{OK: true, Kind: kind, Message: message, Rows: rows}
}

func failure(kind string, err error) Result {
	return Result{OK: false, Kind: kind, Error: err.Error(), Message: "operation failed"}
}

func formatDuty(duty domain.DutyShift) string {
	return strings.Join([]string{duty.DateKey, duty.Shift.Label(), duty.Clinician, string(duty.Status)}, " | ")
}
