package report

import (
	"fmt"
	"sort"

	"hospitalportal/internal/domain"
)

type HealthCheck struct {
	Name    string
	Passed  bool
	Details string
}

type HealthReport struct {
	Checks []HealthCheck
}

func CheckDepartment(department domain.Department, accounts []domain.Account, duties []domain.DutyShift) HealthReport {
	checks := []HealthCheck{
		{Name: "identity", Passed: department.ID != "" && department.Name != "", Details: "department identity"},
		{Name: "status", Passed: department.Status.IsVisible(), Details: "department is visible"},
		{Name: "staffing", Passed: countDepartmentAccounts(department.ID, accounts) > 0, Details: "at least one account"},
		{Name: "schedule", Passed: countPublishedDuties(department.ID, duties) > 0, Details: "at least one published duty"},
	}
	return HealthReport{Checks: checks}
}

func countDepartmentAccounts(id string, accounts []domain.Account) int {
	count := 0
	for _, account := range accounts {
		if account.DepartmentID == id && account.Active {
			count++
		}
	}
	return count
}

func countPublishedDuties(id string, duties []domain.DutyShift) int {
	count := 0
	for _, duty := range duties {
		if duty.DepartmentID == id && duty.Status == domain.RecordPublished {
			count++
		}
	}
	return count
}

func (r HealthReport) Passed() bool {
	for _, check := range r.Checks {
		if !check.Passed {
			return false
		}
	}
	return true
}

func (r HealthReport) Lines() []string {
	lines := make([]string, 0, len(r.Checks)+1)
	state := "unhealthy"
	if r.Passed() {
		state = "healthy"
	}
	lines = append(lines, "overall="+state)
	for _, check := range r.Checks {
		mark := "fail"
		if check.Passed {
			mark = "pass"
		}
		lines = append(lines, fmt.Sprintf("%s=%s (%s)", check.Name, mark, check.Details))
	}
	return lines
}

func RankDepartments(reports map[string]HealthReport) []string {
	ids := make([]string, 0, len(reports))
	for id := range reports {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool {
		left, right := reports[ids[i]], reports[ids[j]]
		if left.Passed() != right.Passed() {
			return left.Passed()
		}
		return ids[i] < ids[j]
	})
	return ids
}

func ExplainHealth(report HealthReport) string {
	if report.Passed() {
		return "all operational checks passed"
	}
	failed := make([]string, 0)
	for _, check := range report.Checks {
		if !check.Passed {
			failed = append(failed, check.Name)
		}
	}
	if len(failed) == 0 {
		return "no checks available"
	}
	return "failed checks: " + fmt.Sprint(failed)
}
