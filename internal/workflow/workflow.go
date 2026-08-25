package workflow

import (
	"hospitalportal/internal/api"
	"hospitalportal/internal/domain"
	"hospitalportal/internal/service"
)

type Runner struct {
	handler *api.Handler
	service *service.Service
}

func NewRunner(s *service.Service) *Runner {
	return &Runner{handler: api.NewHandler(s), service: s}
}

func (r *Runner) DepartmentProfile(request api.DepartmentRequest, lookupID string) []api.Result {
	results := []api.Result{r.handler.CreateDepartment(request)}
	results = append(results, r.handler.FindDepartment(lookupID))
	results = append(results, r.handler.ListDepartmentRows())
	return results
}

func (r *Runner) AccountProvision(request api.AccountRequest) []api.Result {
	created := r.handler.CreateAccount(request)
	results := []api.Result{created}
	if created.OK {
		accounts, err := r.service.ListAccounts(request.DepartmentID)
		if err != nil {
			results = append(results, api.Result{Kind: "account", Message: "account list failed", Error: err.Error()})
		} else {
			rows := make([]string, 0, len(accounts))
			for _, account := range accounts {
				rows = append(rows, account.PermissionSummary())
			}
			results = append(results, api.Result{OK: true, Kind: "account", Message: "accounts listed", Rows: rows})
		}
	}
	return results
}

func (r *Runner) DutyPublishing(request api.DutyRequest) []api.Result {
	created := r.handler.CreateDuty(request)
	results := []api.Result{created}
	if created.OK {
		var duty domain.DutyShift
		duties, err := r.service.ListDutyShifts(request.DepartmentID)
		if err == nil && len(duties) > 0 {
			duty = duties[len(duties)-1]
			results = append(results, r.handler.PublishDuty(duty.ID))
		}
		results = append(results, r.handler.ListDutyRows(request.DepartmentID))
	}
	return results
}

func (r *Runner) AuditRows() ([]string, error) {
	logs, err := r.service.Audit().Recent(0)
	if err != nil {
		return nil, err
	}
	rows := make([]string, 0, len(logs))
	for _, log := range logs {
		rows = append(rows, log.Display())
	}
	return rows, nil
}
