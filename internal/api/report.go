package api

import (
	"fmt"

	"hospitalportal/internal/report"
)

func (h *Handler) OperationsReport(departmentID string) Result {
	data, err := h.service.OperationsReport(departmentID)
	if err != nil {
		return failure("report", err)
	}
	return success("report", data.Header(), data.Lines())
}

func (h *Handler) DepartmentHealth(departmentID string) Result {
	health, err := h.service.DepartmentHealth(departmentID)
	if err != nil {
		return failure("health", err)
	}
	return success("health", fmt.Sprintf("department %s is %s", departmentID, health), nil)
}

func (h *Handler) AccessDecision(accountID, action string) Result {
	decision, err := h.service.AccessDecision(accountID, action)
	if err != nil {
		return failure("access", err)
	}
	state := "denied"
	if decision.Allowed {
		state = "allowed"
	}
	return success("access", state, []string{decision.Username + ": " + decision.Reason})
}

func (h *Handler) HealthReport(departmentID string) Result {
	data, err := h.service.HealthReport(departmentID)
	if err != nil {
		return failure("health", err)
	}
	return success("health", "department health", data.Lines())
}

func ReportLines(data report.OperationsReport) []string {
	return data.Lines()
}
