package domain

import "testing"

func TestDomainModels(t *testing.T) {
	department := NewDepartment("cardio", "Cardiology", "Heart care")
	if err := department.ValidateShape(); err != nil {
		t.Fatal(err)
	}
	paused, err := department.Pause()
	if err != nil || paused.Status != DepartmentPaused {
		t.Fatalf("pause result: %#v %v", paused, err)
	}
	account := NewAccount("acct-1", "cardio", "lee", RoleDoctor)
	if account.PermissionSummary() != "lee/Doctor/active" {
		t.Fatalf("unexpected account summary %q", account.PermissionSummary())
	}
	duty := NewDutyShift("shift-1", "cardio", "2026-08-25", "Lee", ShiftNight)
	if duty.SlotKey() != "cardio|2026-08-25|night" {
		t.Fatalf("unexpected slot key %q", duty.SlotKey())
	}
	log := NewErrorLog("err-1", "code", "message", "context")
	if log.Display() == "" {
		t.Fatal("expected error display")
	}
}
