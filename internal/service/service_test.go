package service

import (
	"path/filepath"
	"testing"

	"hospitalportal/internal/domain"
	"hospitalportal/internal/store"
)

func TestServiceStatusTransitions(t *testing.T) {
	database, err := store.Open(filepath.Join(t.TempDir(), "service.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	serviceLayer := New(database)
	if _, err := serviceLayer.CreateDepartment("cardio", "Cardiology", "Heart care"); err != nil {
		t.Fatal(err)
	}
	if _, err := serviceLayer.ChangeDepartmentStatus("cardio", domain.DepartmentPaused); err != nil {
		t.Fatal(err)
	}
	if _, err := serviceLayer.CreateAccount("acct", "cardio", "lee", domain.RoleDoctor); err == nil {
		t.Fatal("paused department accepted account")
	}
	if _, err := serviceLayer.ChangeDepartmentStatus("cardio", domain.DepartmentActive); err != nil {
		t.Fatal(err)
	}
	if _, err := serviceLayer.CreateAccount("acct", "cardio", "lee", domain.RoleDoctor); err != nil {
		t.Fatal(err)
	}
}
