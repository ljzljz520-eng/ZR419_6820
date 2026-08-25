package report

import (
	"fmt"
	"sort"

	"hospitalportal/internal/domain"
)

type TrendPoint struct {
	Label string
	Value int
}

type PortalMetrics struct {
	Departments int
	Accounts    int
	Active      int
	Published   int
	Errors      int
	Health      map[string]string
}

func CalculateMetrics(departments []domain.Department, accounts []domain.Account, duties []domain.DutyShift, logs []domain.ErrorLog) PortalMetrics {
	metrics := PortalMetrics{Departments: len(departments), Accounts: len(accounts), Health: make(map[string]string)}
	for _, account := range accounts {
		if account.Active {
			metrics.Active++
		}
	}
	for _, duty := range duties {
		if duty.Status == domain.RecordPublished {
			metrics.Published++
		}
	}
	metrics.Errors = len(logs)
	for _, department := range departments {
		metrics.Health[department.ID] = department.Status.String()
	}
	return metrics
}

func (m PortalMetrics) Summary() string {
	return fmt.Sprintf("departments=%d accounts=%d active=%d published=%d errors=%d", m.Departments, m.Accounts, m.Active, m.Published, m.Errors)
}

func (m PortalMetrics) HealthLines() []string {
	keys := make([]string, 0, len(m.Health))
	for key := range m.Health {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	lines := make([]string, 0, len(keys))
	for _, key := range keys {
		lines = append(lines, key+"="+m.Health[key])
	}
	return lines
}

func TrendFromDuties(duties []domain.DutyShift) []TrendPoint {
	counts := make(map[string]int)
	for _, duty := range duties {
		if duty.Status == domain.RecordPublished {
			counts[duty.DateKey]++
		}
	}
	labels := make([]string, 0, len(counts))
	for label := range counts {
		labels = append(labels, label)
	}
	sort.Strings(labels)
	trend := make([]TrendPoint, 0, len(labels))
	for _, label := range labels {
		trend = append(trend, TrendPoint{Label: label, Value: counts[label]})
	}
	return trend
}
