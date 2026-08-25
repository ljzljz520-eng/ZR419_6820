package workflow

import (
	"hospitalportal/internal/api"
)

func (r *Runner) DepartmentReview(departmentID string) []api.Result {
	results := []api.Result{r.handler.OperationsReport(departmentID), r.handler.DepartmentHealth(departmentID)}
	rows, err := r.AuditRows()
	if err != nil {
		results = append(results, api.Result{Kind: "audit", Message: "audit unavailable", Error: err.Error()})
	} else {
		results = append(results, api.Result{OK: true, Kind: "audit", Message: "audit loaded", Rows: rows})
	}
	return results
}
