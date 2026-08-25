package validate

import (
	"testing"

	"hospitalportal/internal/domain"
)

func TestValidation(t *testing.T) {
	if err := DepartmentInput("cardio", "Cardiology", "care"); err != nil {
		t.Fatal(err)
	}
	if err := AccountInput("acct", "cardio", "Doctor", domain.RoleDoctor); err != nil {
		t.Fatal(err)
	}
	if err := DutyInput("shift", "cardio", "2026-08-25", "Lee", domain.ShiftNight); err != nil {
		t.Fatal(err)
	}
	if IsRoleSupported(domain.AccountRole("visitor")) {
		t.Fatal("unknown role accepted")
	}
	candidate := domain.NewDutyShift("two", "cardio", "2026-08-25", "Kim", domain.ShiftNight)
	existing := domain.NewDutyShift("one", "cardio", "2026-08-25", "Lee", domain.ShiftNight)
	if !HasDutyConflict([]domain.DutyShift{existing}, candidate) {
		t.Fatal("expected duty conflict")
	}
}
