package domain

import (
	"sort"
	"strings"
)

type AuditSummary struct {
	Total       int
	ByCode      map[string]int
	ByContext   map[string]int
	LastCode    string
	LastMessage string
}

func SummarizeErrors(logs []ErrorLog) AuditSummary {
	summary := AuditSummary{ByCode: make(map[string]int), ByContext: make(map[string]int)}
	ordered := append([]ErrorLog(nil), logs...)
	sort.Slice(ordered, func(i, j int) bool { return ordered[i].ID < ordered[j].ID })
	for _, log := range ordered {
		summary.Total++
		summary.ByCode[log.Code]++
		context := strings.TrimSpace(log.Context)
		if context == "" {
			context = "global"
		}
		summary.ByContext[context]++
		summary.LastCode = log.Code
		summary.LastMessage = log.Message
	}
	return summary
}

func (s AuditSummary) HasRepeatedCode(threshold int) bool {
	if threshold < 1 {
		threshold = 1
	}
	for _, count := range s.ByCode {
		if count >= threshold {
			return true
		}
	}
	return false
}

func (s AuditSummary) ContextCount(context string) int {
	return s.ByContext[strings.TrimSpace(context)]
}
