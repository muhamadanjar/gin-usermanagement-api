package usecase

import (
	"testing"

	"usermanagement-api/domain/entities"
	"usermanagement-api/domain/repositories"

	"github.com/google/uuid"
)

// Minimal fakes for the ports the authorizer depends on.

type nopUserRepo struct{}

func (nopUserRepo) Create(*entities.User) error                       { return nil }
func (nopUserRepo) FindByID(uuid.UUID) (*entities.User, error)        { return nil, nil }
func (nopUserRepo) FindByEmail(string) (*entities.User, error)        { return nil, nil }
func (nopUserRepo) FindByUsername(string) (*entities.User, error)     { return nil, nil }
func (nopUserRepo) FindAll(int, int) ([]*entities.User, int64, error) { return nil, 0, nil }
func (nopUserRepo) Update(*entities.User) error                       { return nil }
func (nopUserRepo) Delete(uuid.UUID) error                            { return nil }
func (nopUserRepo) AssignRoles(uuid.UUID, []uuid.UUID) error          { return nil }
func (nopUserRepo) AppendRole(uuid.UUID, uuid.UUID) error             { return nil }
func (nopUserRepo) AddTokenHistory(*entities.TokenHistory) error      { return nil }
func (nopUserRepo) FindTokenHistory(uuid.UUID) ([]*entities.TokenHistory, error) {
	return nil, nil
}
func (nopUserRepo) RemoveTokenHistoryByToken(string) error { return nil }

type nopRoleRepo struct{}

func (nopRoleRepo) Create(*entities.Role) error                           { return nil }
func (nopRoleRepo) FindByID(uuid.UUID) (*entities.Role, error)            { return nil, nil }
func (nopRoleRepo) FindByName(string) (*entities.Role, error)             { return nil, nil }
func (nopRoleRepo) FindAll(int, int) ([]*entities.Role, int64, error)     { return nil, 0, nil }
func (nopRoleRepo) Update(*entities.Role) error                           { return nil }
func (nopRoleRepo) Delete(uuid.UUID) error                                { return nil }
func (nopRoleRepo) AssignPermissions(uuid.UUID, []uuid.UUID) error        { return nil }
func (nopRoleRepo) FindRolesByUserID(uuid.UUID) ([]*entities.Role, error) { return nil, nil }
func (nopRoleRepo) FindMembersByRoleID(uuid.UUID) ([]*entities.User, error) {
	return nil, nil
}
func (nopRoleRepo) FindPermissionsByRoleIDs([]uuid.UUID) ([]*entities.Permission, error) {
	return nil, nil
}

func newTestAuthorizer() Authorizer {
	return NewAuthorizer(nopUserRepo{}, nopRoleRepo{}, nil)
}

func TestCheckPermissionSuperuserBypass(t *testing.T) {
	a := newTestAuthorizer()
	principal := &entities.Principal{User: &entities.User{IsSuperuser: true}}
	ok, err := a.CheckPermission(principal, "menu", uuid.New(), "whatever")
	if err != nil || !ok {
		t.Fatalf("superuser should bypass: ok=%v err=%v", ok, err)
	}
}

func TestCheckPermissionRoleGrant(t *testing.T) {
	perm := &entities.Permission{ID: uuid.New(), Name: "users.read"}
	a := newTestAuthorizer()
	principal := &entities.Principal{User: &entities.User{}, Permissions: []*entities.Permission{perm}}
	ok, err := a.CheckPermission(principal, "menu", uuid.New(), "users.read")
	if err != nil || !ok {
		t.Fatalf("role-granted permission should pass: ok=%v err=%v", ok, err)
	}
}

func TestCheckPermissionUnknownDenies(t *testing.T) {
	a := newTestAuthorizer()
	principal := &entities.Principal{User: &entities.User{}}
	ok, err := a.CheckPermission(principal, "menu", uuid.New(), "does.not.exist")
	if err != nil || ok {
		t.Fatalf("unknown permission should deny: ok=%v err=%v", ok, err)
	}
}

func TestPrincipalHasRole(t *testing.T) {
	p := &entities.Principal{Roles: []*entities.Role{{Name: "admin"}}}
	if !p.HasRole("admin") || p.HasRole("user") {
		t.Fatal("HasRole casing/name mismatch")
	}
}

var _ repositories.UserRepository = nopUserRepo{}
var _ repositories.RoleRepository = nopRoleRepo{}
