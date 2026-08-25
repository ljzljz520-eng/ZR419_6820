package service

import (
	"hospitalportal/internal/domain"
	"hospitalportal/internal/report"
)

func (s *Service) OperationsReport(departmentID string) (report.OperationsReport, error) {
	department, err := s.MustGetDepartment(departmentID)
	if err != nil {
		return report.OperationsReport{}, s.recordFailure("report.department", err, departmentID)
	}
	accounts, err := s.store.ListAccounts(departmentID)
	if err != nil {
		return report.OperationsReport{}, s.recordFailure("report.accounts", err, departmentID)
	}
	duties, err := s.store.ListDutyShifts(departmentID)
	if err != nil {
		return report.OperationsReport{}, s.recordFailure("report.duties", err, departmentID)
	}
	logs, err := s.store.FindErrorLogsByContext(departmentID)
	if err != nil {
		return report.OperationsReport{}, s.recordFailure("report.errors", err, departmentID)
	}
	return report.Build(department, accounts, duties, logs), nil
}

func (s *Service) DepartmentHealth(departmentID string) (string, error) {
	reportData, err := s.OperationsReport(departmentID)
	if err != nil {
		return "", err
	}
	return reportData.Snapshot.Health(), nil
}

func (s *Service) AccessDecision(accountID, action string) (domain.AccessDecision, error) {
	account, err := s.GetAccount(accountID)
	if err != nil {
		return domain.AccessDecision{}, err
	}
	department, err := s.MustGetDepartment(account.DepartmentID)
	if err != nil {
		return domain.AccessDecision{}, err
	}
	return account.DecideAccess(action, department), nil
}

func (s *Service) HealthReport(departmentID string) (report.HealthReport, error) {
	department, err := s.MustGetDepartment(departmentID)
	if err != nil {
		return report.HealthReport{}, err
	}
	accounts, err := s.store.ListAccounts(departmentID)
	if err != nil {
		return report.HealthReport{}, s.recordFailure("health.accounts", err, departmentID)
	}
	duties, err := s.store.ListDutyShifts(departmentID)
	if err != nil {
		return report.HealthReport{}, s.recordFailure("health.duties", err, departmentID)
	}
	return report.CheckDepartment(department, accounts, duties), nil
}

func (s *Service) HealthLines(departmentID string) ([]string, error) {
	data, err := s.HealthReport(departmentID)
	if err != nil {
		return nil, err
	}
	return data.Lines(), nil
}

func (s *Service) AuditSummary() (domain.AuditSummary, error) {
	logs, err := s.store.ListErrorLogs()
	if err != nil {
		return domain.AuditSummary{}, s.recordFailure("audit.summary", err, "global")
	}
	return domain.SummarizeErrors(logs), nil
}

func (s *Service) AuditNeedsReview(threshold int) (bool, error) {
	summary, err := s.AuditSummary()
	if err != nil {
		return false, err
	}
	return summary.HasRepeatedCode(threshold), nil
}
