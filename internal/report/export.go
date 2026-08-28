package report

import (
	"fmt"
	"strings"

	"hospitalportal/internal/domain"
)

func ExportDepartment(department domain.Department, snapshot domain.DepartmentSnapshot) string {
	rows := []string{
		"department_id=" + department.ID,
		"department_name=" + department.Name,
		"department_status=" + string(department.Status),
		fmt.Sprintf("account_count=%d", snapshot.AccountCount),
		fmt.Sprintf("duty_count=%d", snapshot.DutyCount),
		fmt.Sprintf("error_count=%d", snapshot.ErrorCount),
		"health=" + snapshot.Health(),
	}
	return strings.Join(rows, "\n")
}

func ExportAccounts(accounts []domain.Account) []string {
	ordered := domain.SortAccountsByRole(accounts)
	rows := make([]string, 0, len(ordered))
	for _, account := range ordered {
		state := "inactive"
		if account.Active {
			state = "active"
		}
		rows = append(rows, fmt.Sprintf("%s|%s|%s|%s", account.ID, account.Username, account.Role, state))
	}
	return rows
}

func ExportDuties(duties []domain.DutyShift) []string {
	rows := make([]string, 0, len(duties))
	for _, duty := range duties {
		rows = append(rows, fmt.Sprintf("%s|%s|%s|%s|%s", duty.ID, duty.DateKey, duty.Shift, duty.Clinician, duty.Status))
	}
	return rows
}

func ExportWarnings(warnings []string) string {
	if len(warnings) == 0 {
		return "warnings=none"
	}
	rows := make([]string, 0, len(warnings)+1)
	rows = append(rows, fmt.Sprintf("warnings=%d", len(warnings)))
	for _, warning := range warnings {
		rows = append(rows, "warning="+warning)
	}
	return strings.Join(rows, "\n")
}

func ExportMetrics(metrics PortalMetrics) []string {
	lines := []string{metrics.Summary()}
	for _, health := range metrics.HealthLines() {
		lines = append(lines, health)
	}
	return lines
}
