package domain

type DepartmentStatus string

const (
	DepartmentActive   DepartmentStatus = "active"
	DepartmentPaused   DepartmentStatus = "paused"
	DepartmentArchived DepartmentStatus = "archived"
)

type AccountRole string

const (
	RoleDoctor        AccountRole = "doctor"
	RoleNurse         AccountRole = "nurse"
	RoleAdministrator AccountRole = "administrator"
)

type ShiftName string

const (
	ShiftMorning ShiftName = "morning"
	ShiftEvening ShiftName = "evening"
	ShiftNight   ShiftName = "night"
)

type RecordStatus string

const (
	RecordDraft     RecordStatus = "draft"
	RecordPublished RecordStatus = "published"
	RecordCancelled RecordStatus = "cancelled"
)

func (s DepartmentStatus) IsVisible() bool {
	return s == DepartmentActive || s == DepartmentPaused
}

func (s DepartmentStatus) String() string {
	return string(s)
}

func (s DepartmentStatus) CanAcceptAccounts() bool {
	return s == DepartmentActive
}

func (r AccountRole) Label() string {
	switch r {
	case RoleDoctor:
		return "Doctor"
	case RoleNurse:
		return "Nurse"
	case RoleAdministrator:
		return "Administrator"
	default:
		return "Unknown"
	}
}

func (s ShiftName) Label() string {
	switch s {
	case ShiftMorning:
		return "Morning"
	case ShiftEvening:
		return "Evening"
	case ShiftNight:
		return "Night"
	default:
		return "Unknown"
	}
}

func (s RecordStatus) IsFinal() bool {
	return s == RecordPublished || s == RecordCancelled
}
