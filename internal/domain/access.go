package domain

import "strings"

type AccessDecision struct {
	Username string
	Role     AccountRole
	Allowed  bool
	Reason   string
}

func (a Account) DecideAccess(action string, department Department) AccessDecision {
	action = strings.ToLower(strings.TrimSpace(action))
	decision := AccessDecision{Username: a.Username, Role: a.Role}
	if !a.Active {
		decision.Reason = "account is inactive"
		return decision
	}
	if department.Status == DepartmentArchived {
		decision.Reason = "department is archived"
		return decision
	}
	switch action {
	case "view", "lookup":
		decision.Allowed = true
		decision.Reason = "read access"
	case "manage_accounts":
		decision.Allowed = a.Role == RoleAdministrator
		decision.Reason = roleReason(decision.Allowed, "administrator role required")
	case "publish_duty":
		decision.Allowed = a.Role == RoleAdministrator || a.Role == RoleDoctor
		decision.Reason = roleReason(decision.Allowed, "clinical coordinator role required")
	default:
		decision.Reason = "unknown action"
	}
	return decision
}

func roleReason(allowed bool, denied string) string {
	if allowed {
		return "role permits action"
	}
	return denied
}

func (a Account) CanView(department Department) bool {
	return a.DecideAccess("view", department).Allowed
}

func (a Account) CanManageAccounts(department Department) bool {
	return a.DecideAccess("manage_accounts", department).Allowed
}

func (a Account) CanPublishDuty(department Department) bool {
	return a.DecideAccess("publish_duty", department).Allowed
}
