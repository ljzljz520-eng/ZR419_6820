package hospitalportal_test

import (
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"hospitalportal/internal/api"
	"hospitalportal/internal/domain"
	"hospitalportal/internal/service"
	"hospitalportal/internal/store"
	"hospitalportal/internal/workflow"
)

func TestWorkflowDepartmentLookup(t *testing.T) {
	path := filepath.Join(t.TempDir(), "department.db")
	database, err := store.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	runner := workflow.NewRunner(service.New(database))
	results := runner.DepartmentProfile(api.DepartmentRequest{ID: "cardio", Name: "Cardiology", Description: "Heart care"}, "cardio")
	if len(results) != 3 || !results[0].OK || !results[1].OK || !strings.Contains(results[1].Rows[0], "cardio") {
		t.Fatalf("department workflow failed: %#v", results)
	}
	if err := database.Close(); err != nil {
		t.Fatal(err)
	}
	reopened, err := store.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer reopened.Close()
	department, found, err := reopened.GetDepartment("cardio")
	if err != nil || !found || department.Name != "Cardiology" {
		t.Fatalf("department persistence failed: %#v %t %v", department, found, err)
	}
}

func TestDepartmentLookupKeepsNotFoundContext(t *testing.T) {
	database, err := store.Open(filepath.Join(t.TempDir(), "missing.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	handler := api.NewHandler(service.New(database))
	result := handler.FindDepartment("neuro-missing")
	if result.OK {
		t.Fatalf("missing department unexpectedly succeeded: %#v", result)
	}
	if !strings.Contains(result.Error, "neuro-missing") {
		t.Fatalf("missing department lost context: %#v", result)
	}
}

func TestWorkflowAccountProvisioning(t *testing.T) {
	database, err := store.Open(filepath.Join(t.TempDir(), "accounts.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	serviceLayer := service.New(database)
	if _, err := serviceLayer.CreateDepartment("cardio", "Cardiology", "Heart care"); err != nil {
		t.Fatal(err)
	}
	runner := workflow.NewRunner(serviceLayer)
	results := runner.AccountProvision(api.AccountRequest{ID: "acct-admin", DepartmentID: "cardio", Username: "chief", Role: domain.RoleAdministrator})
	if len(results) != 2 || !results[0].OK || !results[1].OK {
		t.Fatalf("account workflow failed: %#v", results)
	}
	if !strings.Contains(results[1].Rows[0], "Administrator") {
		t.Fatalf("account role not displayed: %#v", results)
	}
}

func TestWorkflowDutyPublishing(t *testing.T) {
	database, err := store.Open(filepath.Join(t.TempDir(), "duties.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	serviceLayer := service.New(database)
	if _, err := serviceLayer.CreateDepartment("cardio", "Cardiology", "Heart care"); err != nil {
		t.Fatal(err)
	}
	runner := workflow.NewRunner(serviceLayer)
	results := runner.DutyPublishing(api.DutyRequest{ID: "shift-1", DepartmentID: "cardio", DateKey: "2026-08-25", Clinician: "Lee", Shift: domain.ShiftMorning})
	if len(results) != 3 || !results[0].OK || !results[1].OK || !results[2].OK {
		t.Fatalf("duty workflow failed: %#v", results)
	}
	if _, err := serviceLayer.CreateDutyShift("shift-2", "cardio", "2026-08-25", "Kim", domain.ShiftMorning); err == nil {
		t.Fatal("duty conflict was not enforced")
	} else if !errors.Is(err, service.ErrDutyConflict{Slot: "cardio|2026-08-25|morning"}) {
		t.Fatalf("unexpected conflict error: %v", err)
	}
}
