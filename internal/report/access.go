package report

import (
	"fmt"
	"sort"

	"hospitalportal/internal/domain"
)

type PermissionMatrix struct {
	Actions []string
	Rows    []PermissionRow
}

type PermissionRow struct {
	Username string
	Role     domain.AccountRole
	Allowed  map[string]bool
}

func BuildPermissionMatrix(department domain.Department, accounts []domain.Account) PermissionMatrix {
	actions := []string{"view", "lookup", "manage_accounts", "publish_duty"}
	rows := make([]PermissionRow, 0)
	for _, account := range domain.SortAccountsByRole(accounts) {
		if account.DepartmentID != department.ID {
			continue
		}
		allowed := make(map[string]bool, len(actions))
		for _, action := range actions {
			allowed[action] = account.DecideAccess(action, department).Allowed
		}
		rows = append(rows, PermissionRow{Username: account.Username, Role: account.Role, Allowed: allowed})
	}
	return PermissionMatrix{Actions: actions, Rows: rows}
}

func (m PermissionMatrix) Lines() []string {
	lines := make([]string, 0, len(m.Rows)+1)
	lines = append(lines, fmt.Sprintf("permissions: %s", m.Actions))
	for _, row := range m.Rows {
		values := make([]string, 0, len(m.Actions))
		for _, action := range m.Actions {
			state := "deny"
			if row.Allowed[action] {
				state = "allow"
			}
			values = append(values, action+"="+state)
		}
		lines = append(lines, row.Username+"/"+row.Role.Label()+" "+fmt.Sprint(values))
	}
	return lines
}

func SortCoverageByOpenShifts(coverage []Coverage) []Coverage {
	copyOf := append([]Coverage(nil), coverage...)
	sort.SliceStable(copyOf, func(i, j int) bool {
		if len(copyOf[i].OpenShifts) == len(copyOf[j].OpenShifts) {
			return copyOf[i].DateKey < copyOf[j].DateKey
		}
		return len(copyOf[i].OpenShifts) > len(copyOf[j].OpenShifts)
	})
	return copyOf
}
