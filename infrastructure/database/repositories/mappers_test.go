package repositories

import (
	"testing"
	"time"

	"usermanagement-api/infrastructure/database/models"

	"github.com/google/uuid"
)

func TestUserRoundTrip(t *testing.T) {
	now := time.Now()
	roleID := uuid.New()
	m := &models.UserModel{
		ID:          uuid.New(),
		Username:    "alice",
		Email:       "alice@example.com",
		Password:    "hash",
		FirstName:   "A",
		LastName:    "Z",
		IsSuperuser: true,
		IsActive:    true,
		CreatedAt:   now,
		UpdatedAt:   now,
		Roles:       []*models.RoleModel{{ID: roleID, Name: "admin", CreatedAt: now, UpdatedAt: now}},
	}

	u := toEntityUser(m)
	if u.ID != m.ID || u.Username != "alice" || !u.IsSuperuser || !u.IsActive {
		t.Fatalf("entity mapping mismatch: %+v", u)
	}
	if len(u.Roles) != 1 || u.Roles[0].ID != roleID {
		t.Fatalf("roles not mapped: %+v", u.Roles)
	}

	back := toModelUser(u)
	if back.ID != m.ID || back.Email != m.Email || len(back.Roles) != 1 {
		t.Fatalf("model roundtrip broke: %+v", back)
	}
}

func TestMenuChildrenMapped(t *testing.T) {
	parentID := uuid.New()
	m := &models.MenuModel{
		ID:        uuid.New(),
		Name:      "parent",
		IsActive:  true,
		IsVisible: true,
		Children:  []*models.MenuModel{{ID: uuid.New(), Name: "child", ParentID: &parentID}},
	}

	e := toEntityMenu(m)
	if len(e.Children) != 1 || e.Children[0].Name != "child" {
		t.Fatalf("children not mapped: %+v", e.Children)
	}
	if !e.IsActive || !e.IsVisible {
		t.Fatalf("flags not mapped: %+v", e)
	}
}