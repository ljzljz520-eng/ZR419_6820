package domain

import (
	"fmt"
	"sort"
	"strings"
)

type DepartmentSnapshot struct {
	DepartmentID string
	Status       DepartmentStatus
	AccountCount int
	DutyCount    int
	ErrorCount   int
}

func (d Department) Snapshot(accounts []Account, duties []DutyShift, logs []ErrorLog) DepartmentSnapshot {
	snapshot := DepartmentSnapshot{DepartmentID: d.ID, Status: d.Status}
	for _, account := range accounts {
		if account.DepartmentID == d.ID {
			snapshot.AccountCount++
		}
	}
	for _, duty := range duties {
		if duty.DepartmentID == d.ID && duty.Status != RecordCancelled {
			snapshot.DutyCount++
		}
	}
	for _, log := range logs {
		if strings.Contains(log.Context, d.ID) {
			snapshot.ErrorCount++
		}
	}
	return snapshot
}

func (s DepartmentSnapshot) Health() string {
	if s.Status == DepartmentArchived {
		return "archived"
	}
	if s.AccountCount == 0 {
		return "unstaffed"
	}
	if s.DutyCount == 0 {
		return "uncovered"
	}
	if s.ErrorCount > 3 {
		return "attention"
	}
	return "ready"
}

func (s DepartmentSnapshot) String() string {
	return fmt.Sprintf("%s %s accounts=%d duties=%d errors=%d", s.DepartmentID, s.Health(), s.AccountCount, s.DutyCount, s.ErrorCount)
}

func SortAccountsByRole(accounts []Account) []Account {
	copyOf := append([]Account(nil), accounts...)
	sort.SliceStable(copyOf, func(i, j int) bool {
		if copyOf[i].Role == copyOf[j].Role {
			return copyOf[i].Username < copyOf[j].Username
		}
		return roleRank(copyOf[i].Role) < roleRank(copyOf[j].Role)
	})
	return copyOf
}

func roleRank(role AccountRole) int {
	switch role {
	case RoleAdministrator:
		return 0
	case RoleDoctor:
		return 1
	case RoleNurse:
		return 2
	default:
		return 3
	}
}

func GroupDutiesByDate(duties []DutyShift) map[string][]DutyShift {
	grouped := make(map[string][]DutyShift)
	for _, duty := range duties {
		grouped[duty.DateKey] = append(grouped[duty.DateKey], duty)
	}
	for date, items := range grouped {
		sort.Slice(items, func(i, j int) bool { return items[i].Shift < items[j].Shift })
		grouped[date] = items
	}
	return grouped
}
