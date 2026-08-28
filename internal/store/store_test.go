package store

import (
	"path/filepath"
	"testing"

	"hospitalportal/internal/domain"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "portal.db")
	database, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := database.SaveDepartment(domain.NewDepartment("cardio", "Cardiology", "Heart care")); err != nil {
		t.Fatal(err)
	}
	if err := database.SaveAccount(domain.NewAccount("acct", "cardio", "lee", domain.RoleDoctor)); err != nil {
		t.Fatal(err)
	}
	if err := database.SaveDutyShift(domain.NewDutyShift("shift", "cardio", "2026-08-25", "Lee", domain.ShiftMorning)); err != nil {
		t.Fatal(err)
	}
	if err := database.SaveErrorLog(domain.NewErrorLog("err", "sample", "message", "cardio")); err != nil {
		t.Fatal(err)
	}
	if err := database.Close(); err != nil {
		t.Fatal(err)
	}
	reopened, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer reopened.Close()
	department, found, err := reopened.GetDepartment("cardio")
	if err != nil || !found || department.Name != "Cardiology" {
		t.Fatalf("department did not survive reopen: %#v %v %t", department, err, found)
	}
	accounts, err := reopened.ListAccounts("cardio")
	if err != nil || len(accounts) != 1 {
		t.Fatalf("accounts did not survive reopen: %v", err)
	}
	duties, err := reopened.ListDutyShifts("cardio")
	if err != nil || len(duties) != 1 {
		t.Fatalf("duties did not survive reopen: %v", err)
	}
	logs, err := reopened.ListErrorLogs()
	if err != nil || len(logs) != 1 {
		t.Fatalf("logs did not survive reopen: %v", err)
	}
}
