package entities

import "strings"

// Principal is the authenticated subject: the user together with the roles
// and permissions resolved from their active roles. It is produced by the
// authorizer use case and consumed by authorization middleware — no database
// access happens inside handlers.
type Principal struct {
	User        *User
	Roles       []*Role
	Permissions []*Permission
}

func (p *Principal) IsSuperuser() bool {
	return p != nil && p.User != nil && p.User.IsSuperuser
}

// HasRole reports whether the principal holds at least one of the given role
// names (case-insensitive).
func (p *Principal) HasRole(names ...string) bool {
	if p == nil {
		return false
	}
	for _, role := range p.Roles {
		for _, name := range names {
			if strings.EqualFold(role.Name, name) {
				return true
			}
		}
	}
	return false
}

// HasPermission reports whether any of the principal's roles carries the
// named permission.
func (p *Principal) HasPermission(name string) bool {
	if p == nil {
		return false
	}
	for _, perm := range p.Permissions {
		if perm.Name == name {
			return true
		}
	}
	return false
}
