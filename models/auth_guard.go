package models

import (
	"database/sql"
	"strconv"
)

type AuthGuard struct {
	SameUserRole     string
	SameCompanyRoles []string
	OverridingRoles  []string
}

type AuthCheck struct {
	AccessorRole string
	AccessorID   string
	OwnerRole    string
	OwnerID      string
}

func (ag *AuthGuard) CanAccess(db *sql.DB, ac AuthCheck) bool {
	if ac.AccessorRole == ag.SameUserRole && ac.AccessorID == ac.OwnerID {
		// Owner is accessing
		return true
	}

	if inArray(ac.AccessorRole, ag.OverridingRoles) {
		return true
	}

	if !inArray(ac.AccessorRole, ag.SameCompanyRoles) {
		return false // We won't bother try the last checks if its a different company
	} else if ac.OwnerRole == "company" { // Special case for entities accessible by users from the same company
		return strconv.Itoa(GetCompanyIDFromID(db, ac.AccessorID, ac.AccessorRole)) == ac.OwnerID
	}

	return GetCompanyIDFromID(db, ac.AccessorID, ac.AccessorRole) == GetCompanyIDFromID(db, ac.OwnerID, ac.OwnerRole)
}

func inArray(needle string, haystack []string) bool {
	for _, potentialNeedle := range haystack {
		if potentialNeedle == needle {
			return true
		}
	}
	return false
}

func (ag AuthGuard) AuthInfo() string {
	req := "You need to meet the following auth: "
	if ag.SameUserRole != "" {
		req += "same user of type " + ag.SameUserRole
	}
	if len(ag.SameCompanyRoles) > 0 {
		req += "same company with user type/s "
		for _, role := range ag.SameCompanyRoles {
			req += " " + role
		}
	}
	if len(ag.OverridingRoles) > 0 {
		req += "or user type/s "
		for _, role := range ag.OverridingRoles {
			req += " " + role
		}
	}
	return req + "."
}
