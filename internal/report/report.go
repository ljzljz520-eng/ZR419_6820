package report

import (
	"fmt"
	"sort"
	"strings"

	"hospitalportal/internal/domain"
)

type Coverage struct {
	DateKey    string
	Department string
	Published  int
	Draft      int
	Cancelled  int
	OpenShifts []domain.ShiftName
}

type AccessSummary struct {
	DepartmentID   string
	Total          int
	Active         int
	Doctors        int
	Nurses         int
	Administrators int
}

type OperationsReport struct {
	Department domain.Department
	Snapshot   domain.DepartmentSnapshot
	Coverage   []Coverage
	Access     AccessSummary
	Warnings   []string
}

func Build(department domain.Department, accounts []domain.Account, duties []domain.DutyShift, logs []domain.ErrorLog) OperationsReport {
	result := OperationsReport{Department: department, Snapshot: department.Snapshot(accounts, duties, logs)}
	result.Access = summarizeAccess(department.ID, accounts)
	result.Coverage = buildCoverage(department.ID, duties)
	result.Warnings = warnings(result.Snapshot, result.Coverage)
	return result
}

func summarizeAccess(departmentID string, accounts []domain.Account) AccessSummary {
	result := AccessSummary{DepartmentID: departmentID}
	for _, account := range accounts {
		if account.DepartmentID != departmentID {
			continue
		}
		result.Total++
		if account.Active {
			result.Active++
		}
		switch account.Role {
		case domain.RoleDoctor:
			result.Doctors++
		case domain.RoleNurse:
			result.Nurses++
		case domain.RoleAdministrator:
			result.Administrators++
		}
	}
	return result
}

func buildCoverage(departmentID string, duties []domain.DutyShift) []Coverage {
	grouped := domain.GroupDutiesByDate(duties)
	dates := make([]string, 0, len(grouped))
	for date := range grouped {
		dates = append(dates, date)
	}
	sort.Strings(dates)
	coverage := make([]Coverage, 0, len(dates))
	for _, date := range dates {
		item := Coverage{DateKey: date, Department: departmentID}
		seen := make(map[domain.ShiftName]bool)
		for _, duty := range grouped[date] {
			seen[duty.Shift] = true
			switch duty.Status {
			case domain.RecordPublished:
				item.Published++
			case domain.RecordDraft:
				item.Draft++
			case domain.RecordCancelled:
				item.Cancelled++
			}
		}
		for _, shift := range []domain.ShiftName{domain.ShiftMorning, domain.ShiftEvening, domain.ShiftNight} {
			if !seen[shift] {
				item.OpenShifts = append(item.OpenShifts, shift)
			}
		}
		coverage = append(coverage, item)
	}
	return coverage
}

func warnings(snapshot domain.DepartmentSnapshot, coverage []Coverage) []string {
	warnings := make([]string, 0)
	if snapshot.Status == domain.DepartmentPaused {
		warnings = append(warnings, "department is paused")
	}
	if snapshot.AccountCount == 0 {
		warnings = append(warnings, "no accounts assigned")
	}
	if len(coverage) == 0 {
		warnings = append(warnings, "no duty schedule exists")
	}
	for _, item := range coverage {
		if len(item.OpenShifts) > 0 {
			warnings = append(warnings, fmt.Sprintf("%s has %d open shifts", item.DateKey, len(item.OpenShifts)))
		}
	}
	if snapshot.ErrorCount > 0 {
		warnings = append(warnings, fmt.Sprintf("%d logged errors need review", snapshot.ErrorCount))
	}
	return warnings
}

func (r OperationsReport) Header() string {
	return fmt.Sprintf("%s | %s | %s", r.Department.ID, r.Department.Name, r.Snapshot.Health())
}

func (r OperationsReport) Lines() []string {
	lines := []string{r.Header(), fmt.Sprintf("accounts: %d total, %d active", r.Access.Total, r.Access.Active), fmt.Sprintf("roles: doctors=%d nurses=%d administrators=%d", r.Access.Doctors, r.Access.Nurses, r.Access.Administrators)}
	for _, coverage := range r.Coverage {
		lines = append(lines, fmt.Sprintf("%s: published=%d draft=%d cancelled=%d open=%s", coverage.DateKey, coverage.Published, coverage.Draft, coverage.Cancelled, shiftNames(coverage.OpenShifts)))
	}
	for _, warning := range r.Warnings {
		lines = append(lines, "warning: "+warning)
	}
	return lines
}

func shiftNames(shifts []domain.ShiftName) string {
	if len(shifts) == 0 {
		return "none"
	}
	parts := make([]string, 0, len(shifts))
	for _, shift := range shifts {
		parts = append(parts, string(shift))
	}
	return strings.Join(parts, ",")
}
