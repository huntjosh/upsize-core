package tests

import (
	"testing"
	"upsizeAPI/models"
)

func TestAuthGuardSameUser(t *testing.T) {
	ag := models.AuthGuard{SameUserRole: "contractor", SameCompanyRoles: []string{"manager"}, OverridingRoles: []string{"admin"}}
	if !ag.CanAccess(a.DB, models.AuthCheck{AccessorRole: "contractor", AccessorID: "1", OwnerRole: "contractor", OwnerID: "1"}) {
		t.Errorf("expected access")
	}
}

func TestAuthGuardSameUserTypeDiffUser(t *testing.T) {
	ag := models.AuthGuard{SameUserRole: "contractor", SameCompanyRoles: []string{"manager"}, OverridingRoles: []string{"admin"}}
	if ag.CanAccess(a.DB, models.AuthCheck{AccessorRole: "contractor", AccessorID: "2", OwnerRole: "contractor", OwnerID: "1"}) {
		t.Errorf("expected no access")
	}
}

func TestAuthGuardSameCompany(t *testing.T) {
	ag := models.AuthGuard{SameUserRole: "contractor", SameCompanyRoles: []string{"manager"}, OverridingRoles: []string{"admin"}}
	if !ag.CanAccess(a.DB, models.AuthCheck{AccessorRole: "manager", AccessorID: "1", OwnerRole: "contractor", OwnerID: "1"}) {
		t.Errorf("expected access")
	}
}

func TestAuthGuardOverriding(t *testing.T) {
	ag := models.AuthGuard{SameUserRole: "contractor", SameCompanyRoles: []string{"manager"}, OverridingRoles: []string{"admin"}}
	if !ag.CanAccess(a.DB, models.AuthCheck{AccessorRole: "admin", OwnerRole: "contractor", OwnerID: "1"}) {
		t.Errorf("expected access")
	}
}
