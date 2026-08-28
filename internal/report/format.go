package report

import (
	"fmt"
	"strings"

	"hospitalportal/internal/domain"
)

func FormatAccess(summary AccessSummary) string {
	return fmt.Sprintf("%s accounts=%d active=%d doctors=%d nurses=%d administrators=%d", summary.DepartmentID, summary.Total, summary.Active, summary.Doctors, summary.Nurses, summary.Administrators)
}

func FormatCoverage(coverage Coverage) string {
	return fmt.Sprintf("%s published=%d draft=%d cancelled=%d open=%s", coverage.DateKey, coverage.Published, coverage.Draft, coverage.Cancelled, shiftNames(coverage.OpenShifts))
}

func FormatReport(report OperationsReport) string {
	return strings.Join(report.Lines(), "\n")
}

func BuildDepartmentIndex(departments []domain.Department) map[string]string {
	index := make(map[string]string, len(departments))
	for _, department := range departments {
		index[department.ID] = department.Name
	}
	return index
}

func VisibleDepartments(departments []domain.Department) []domain.Department {
	visible := make([]domain.Department, 0, len(departments))
	for _, department := range departments {
		if department.Status.IsVisible() {
			visible = append(visible, department)
		}
	}
	return visible
}
