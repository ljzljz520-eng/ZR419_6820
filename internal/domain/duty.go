package domain

import (
	"fmt"
	"strings"
)

type DutyShift struct {
	ID           string
	DepartmentID string
	DateKey      string
	Clinician    string
	Shift        ShiftName
	Status       RecordStatus
	Version      int
}

func NewDutyShift(id, departmentID, dateKey, clinician string, shift ShiftName) DutyShift {
	return DutyShift{ID: strings.TrimSpace(id), DepartmentID: strings.TrimSpace(departmentID), DateKey: strings.TrimSpace(dateKey), Clinician: strings.TrimSpace(clinician), Shift: shift, Status: RecordDraft, Version: 1}
}

func (d DutyShift) ValidateShape() error {
	if d.ID == "" || d.DepartmentID == "" || d.DateKey == "" || d.Clinician == "" {
		return fmt.Errorf("duty shift %s has incomplete fields", d.ID)
	}
	if d.Version < 1 {
		return fmt.Errorf("duty shift %s version must be positive", d.ID)
	}
	return nil
}

func (d DutyShift) Publish() (DutyShift, error) {
	if d.Status == RecordCancelled {
		return d, fmt.Errorf("duty shift %s is cancelled", d.ID)
	}
	d.Status = RecordPublished
	d.Version++
	return d, nil
}

func (d DutyShift) Cancel() (DutyShift, error) {
	if d.Status == RecordCancelled {
		return d, fmt.Errorf("duty shift %s is already cancelled", d.ID)
	}
	d.Status = RecordCancelled
	d.Version++
	return d, nil
}

func (d DutyShift) SlotKey() string {
	return d.DepartmentID + "|" + d.DateKey + "|" + string(d.Shift)
}
