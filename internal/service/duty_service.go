package service

import (
	"hospitalportal/internal/domain"
	"hospitalportal/internal/validate"
)

func (s *Service) CreateDutyShift(id, departmentID, dateKey, clinician string, shift domain.ShiftName) (domain.DutyShift, error) {
	if err := validate.DutyInput(id, departmentID, dateKey, clinician, shift); err != nil {
		return domain.DutyShift{}, s.recordFailure("duty.validation", err, id)
	}
	department, err := s.MustGetDepartment(departmentID)
	if err != nil {
		return domain.DutyShift{}, s.recordFailure("duty.department", err, departmentID)
	}
	if err := validate.CanPublishDuty(department, clinician); err != nil {
		return domain.DutyShift{}, s.recordFailure("duty.permission", err, id)
	}
	duty := domain.NewDutyShift(id, departmentID, dateKey, clinician, shift)
	duties, err := s.store.ListDutyShifts(departmentID)
	if err != nil {
		return domain.DutyShift{}, s.recordFailure("duty.list", err, departmentID)
	}
	if validate.HasDutyConflict(duties, duty) {
		err := s.recordFailure("duty.conflict", ErrDutyConflict{Slot: duty.SlotKey()}, id)
		return domain.DutyShift{}, err
	}
	if err := s.store.SaveDutyShift(duty); err != nil {
		return domain.DutyShift{}, s.recordFailure("duty.save", err, id)
	}
	return duty, nil
}

type ErrDutyConflict struct {
	Slot string
}

func (e ErrDutyConflict) Error() string {
	return "duty slot already assigned: " + e.Slot
}

func (s *Service) GetDutyShift(id string) (domain.DutyShift, error) {
	duty, found, err := s.store.GetDutyShift(id)
	if err != nil {
		return domain.DutyShift{}, s.recordFailure("duty.lookup", err, id)
	}
	if !found {
		return domain.DutyShift{}, s.recordFailure("duty.not_found", DutyNotFoundError{ID: id}, id)
	}
	return duty, nil
}

func (s *Service) ListDutyShifts(departmentID string) ([]domain.DutyShift, error) {
	duties, err := s.store.ListDutyShifts(departmentID)
	if err != nil {
		return nil, s.recordFailure("duty.list", err, departmentID)
	}
	return duties, nil
}

func (s *Service) PublishDutyShift(id string) (domain.DutyShift, error) {
	duty, err := s.GetDutyShift(id)
	if err != nil {
		return domain.DutyShift{}, err
	}
	changed, err := duty.Publish()
	if err != nil {
		return domain.DutyShift{}, s.recordFailure("duty.publish", err, id)
	}
	if err := s.store.SaveDutyShift(changed); err != nil {
		return domain.DutyShift{}, s.recordFailure("duty.publish.save", err, id)
	}
	return changed, nil
}

func (s *Service) CancelDutyShift(id string) (domain.DutyShift, error) {
	duty, err := s.GetDutyShift(id)
	if err != nil {
		return domain.DutyShift{}, err
	}
	changed, err := duty.Cancel()
	if err != nil {
		return domain.DutyShift{}, s.recordFailure("duty.cancel", err, id)
	}
	if err := s.store.SaveDutyShift(changed); err != nil {
		return domain.DutyShift{}, s.recordFailure("duty.cancel.save", err, id)
	}
	return changed, nil
}
