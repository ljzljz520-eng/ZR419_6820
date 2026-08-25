package api

import (
	"fmt"
)

func (h *Handler) MaintenanceSummary() Result {
	summary, err := h.service.MaintenanceSummary()
	if err != nil {
		return failure("maintenance", err)
	}
	return success("maintenance", "inventory", []string{summary})
}

func (h *Handler) PermissionMatrix(departmentID string) Result {
	matrix, err := h.service.PermissionMatrix(departmentID)
	if err != nil {
		return failure("permissions", err)
	}
	return success("permissions", fmt.Sprintf("department %s", departmentID), matrix.Lines())
}

func (h *Handler) SeedDefaults() Result {
	if err := h.service.SeedDefaults(); err != nil {
		return failure("maintenance", err)
	}
	return success("maintenance", "reference departments ready", nil)
}
