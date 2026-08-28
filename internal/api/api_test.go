package api

import (
	"path/filepath"
	"strings"
	"testing"

	"hospitalportal/internal/service"
	"hospitalportal/internal/store"
)

func TestAPIPresentation(t *testing.T) {
	database, err := store.Open(filepath.Join(t.TempDir(), "api.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	handler := NewHandler(service.New(database))
	result := handler.CreateDepartment(DepartmentRequest{ID: "cardio", Name: "Cardiology", Description: "Heart care"})
	if !result.OK || !strings.Contains(Render(result), "cardio") {
		t.Fatalf("unexpected presentation: %#v", result)
	}
	if !strings.Contains(RenderErrorLogs(nil), "No error logs") {
		t.Fatal("empty error log rendering changed")
	}
}
