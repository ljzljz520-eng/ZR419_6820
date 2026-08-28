package workflow

import (
	"path/filepath"
	"testing"

	"hospitalportal/internal/api"
	"hospitalportal/internal/domain"
	"hospitalportal/internal/service"
	"hospitalportal/internal/store"
)

func TestWorkflowRunner(t *testing.T) {
	database, err := store.Open(filepath.Join(t.TempDir(), "workflow.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	runner := NewRunner(service.New(database))
	departmentResults := runner.DepartmentProfile(api.DepartmentRequest{ID: "cardio", Name: "Cardiology", Description: "Heart care"}, "cardio")
	if len(departmentResults) != 3 || !departmentResults[0].OK {
		t.Fatalf("department workflow incomplete: %#v", departmentResults)
	}
	accountResults := runner.AccountProvision(api.AccountRequest{ID: "acct", DepartmentID: "cardio", Username: "lee", Role: domain.RoleDoctor})
	if len(accountResults) != 2 || !accountResults[0].OK {
		t.Fatalf("account workflow incomplete: %#v", accountResults)
	}
}
