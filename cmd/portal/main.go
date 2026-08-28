package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"hospitalportal/internal/api"
	"hospitalportal/internal/domain"
	"hospitalportal/internal/service"
	"hospitalportal/internal/store"
	"hospitalportal/internal/workflow"
)

func main() {
	path := flag.String("db", filepath.Join(os.TempDir(), "hospitalportal.db"), "database file")
	flag.Parse()
	database, err := store.Open(*path)
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}
	defer database.Close()
	serviceLayer := service.New(database)
	runner := workflow.NewRunner(serviceLayer)
	for _, result := range runner.DepartmentProfile(api.DepartmentRequest{ID: "cardio", Name: "Cardiology", Description: "Heart care"}, "cardio") {
		fmt.Println(api.Render(result))
	}
	for _, result := range runner.AccountProvision(api.AccountRequest{ID: "acct-admin", DepartmentID: "cardio", Username: "chief", Role: domain.RoleAdministrator}) {
		fmt.Println(api.Render(result))
	}
	for _, result := range runner.DutyPublishing(api.DutyRequest{ID: "shift-001", DepartmentID: "cardio", DateKey: "2026-08-25", Clinician: "Dr Lee", Shift: domain.ShiftMorning}) {
		fmt.Println(api.Render(result))
	}
	if rows, err := runner.AuditRows(); err == nil {
		fmt.Println(api.RenderErrorLogs(rows))
	}
}
