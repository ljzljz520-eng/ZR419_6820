package audit

import (
	"path/filepath"
	"testing"

	"hospitalportal/internal/store"
)

func TestAuditLog(t *testing.T) {
	database, err := store.Open(filepath.Join(t.TempDir(), "audit.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	logger := NewLogger(database)
	if _, err := logger.Record("lookup.not_found", "missing department", "cardio"); err != nil {
		t.Fatal(err)
	}
	found, err := logger.HasCode("lookup.not_found")
	if err != nil || !found {
		t.Fatalf("expected audit code: %v %t", err, found)
	}
}
