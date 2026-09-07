package usecase

import (
	"errors"
	"testing"

	"usermanagement-api/domain/entities"
	"usermanagement-api/domain/repositories"

	"github.com/google/uuid"
)

// Minimal fakes for the ports the authorizer depends on; only permission
// lookup and model-permission check are exercised by CheckPermission.

type fakePermissionRepo struct {
	byName map[string]*entities.Permission
}

func (f *fakePermissionRepo) FindByName(name string) (*entities.Permission, error) {
	if p, ok := f.byName[name]; ok {
		return p, nil
	}
	return nil, errors.New("not found")
}
func (f *fakePermissionRepo) Create(*entities.Permission) error                { return nil }
func (f *fakePermissionRepo) FindByID(uuid.UUID) (*entities.Permission, error) { return nil, nil }
func (f *fakePermissionRepo) FindAll(int, int) ([]*entities.Permission, int64, error) {
	return nil, 0, nil
}
func (f *fakePermissionRepo) Update(*entities.Permission) error { return nil }
func (f *fakePermissionRepo) Delete(uuid.UUID) error            { return nil }

type fakeModelPermRepo struct {
	has bool
}

func (f *fakeModelPermRepo) Create(*entities.ModelPermission) error { return nil }
func (f *fakeModelPermRepo) FindByID(uuid.UUID) (*entities.ModelPermission, error) {
	return nil, nil
}
func (f *fakeModelPermRepo) FindByModelTypeAndModelID(string, uuid.UUID) ([]*entities.ModelPermission, error) {
	return nil, nil
}
func (f *fakeModelPermRepo) FindAll(int, int) ([]*entities.ModelPermission, int64, error) {
	return nil, 0, nil
}
func (f *fakeModelPermRepo) Update(*entities.ModelPermission) error { return nil }
func (f *fakeModelPermRepo) Delete(uuid.UUID) error                 { return nil }
func (f *fakeModelPermRepo) CheckPermission(string, uuid.UUID, uuid.UUID) (bool, error) {
	return f.has, nil
}

type nopUserRepo struct{}

func (nopUserRepo) Create(*entities.User) error                       { return nil }
func (nopUserRepo) FindByID(uuid.UUID) (*entities.User, error)        { return nil, nil }
func (nopUserRepo) FindByEmail(string) (*entities.User, error)        { return nil, nil }
func (nopUserRepo) FindByUsername(string) (*entities.User, error)     { return nil, nil }
func (nopUserRepo) FindAll(int, int) ([]*entities.User, int64, error) { return nil, 0, nil }
func (nopUserRepo) Update(*entities.User) error                       { return nil }
func (nopUserRepo) Delete(uuid.UUID) error                            { return nil }
func (nopUserRepo) AssignRoles(uuid.UUID, []uuid.UUID) error          { return nil }

type nopRoleRepo struct{}

func (nopRoleRepo) Create(*entities.Role) error                           { return nil }
func (nopRoleRepo) FindByID(uuid.UUID) (*entities.Role, error)            { return nil, nil }
func (nopRoleRepo) FindByName(string) (*entities.Role, error)             { return nil, nil }
func (nopRoleRepo) FindAll(int, int) ([]*entities.Role, int64, error)     { return nil, 0, nil }
func (nopRoleRepo) Update(*entities.Role) error                           { return nil }
func (nopRoleRepo) Delete(uuid.UUID) error                                { return nil }
func (nopRoleRepo) AssignPermissions(uuid.UUID, []uuid.UUID) error        { return nil }
func (nopRoleRepo) FindRolesByUserID(uuid.UUID) ([]*entities.Role, error) { return nil, nil }
func (nopRoleRepo) FindPermissionsByRoleIDs([]uuid.UUID) ([]*entities.Permission, error) {
	return nil, nil
}

func newTestAuthorizer(permRepo repositories.PermissionRepository, modelPermRepo repositories.ModelPermissionRepository) Authorizer {
	return NewAuthorizer(nopUserRepo{}, nopRoleRepo{}, permRepo, modelPermRepo, nil)
}

func TestCheckPermissionSuperuserBypass(t *testing.T) {
	a := newTestAuthorizer(&fakePermissionRepo{}, &fakeModelPermRepo{has: false})
	principal := &entities.Principal{User: &entities.User{IsSuperuser: true}}
	ok, err := a.CheckPermission(principal, "menu", uuid.New(), "whatever")
	if err != nil || !ok {
		t.Fatalf("superuser should bypass: ok=%v err=%v", ok, err)
	}
}

func TestCheckPermissionRoleGrant(t *testing.T) {
	perm := &entities.Permission{ID: uuid.New(), Name: "users.read"}
	a := newTestAuthorizer(
		&fakePermissionRepo{byName: map[string]*entities.Permission{"users.read": perm}},
		&fakeModelPermRepo{has: false},
	)
	principal := &entities.Principal{User: &entities.User{}, Permissions: []*entities.Permission{perm}}
	ok, err := a.CheckPermission(principal, "menu", uuid.New(), "users.read")
	if err != nil || !ok {
		t.Fatalf("role-granted permission should pass: ok=%v err=%v", ok, err)
	}
}

func TestCheckPermissionModelFallback(t *testing.T) {
	perm := &entities.Permission{ID: uuid.New(), Name: "menu.read"}
	a := newTestAuthorizer(
		&fakePermissionRepo{byName: map[string]*entities.Permission{"menu.read": perm}},
		&fakeModelPermRepo{has: true},
	)
	principal := &entities.Principal{User: &entities.User{}}
	ok, err := a.CheckPermission(principal, "menu", uuid.New(), "menu.read")
	if err != nil || !ok {
		t.Fatalf("model permission should grant: ok=%v err=%v", ok, err)
	}
}

func TestCheckPermissionUnknownDenies(t *testing.T) {
	a := newTestAuthorizer(&fakePermissionRepo{}, &fakeModelPermRepo{has: false})
	principal := &entities.Principal{User: &entities.User{}}
	ok, err := a.CheckPermission(principal, "menu", uuid.New(), "does.not.exist")
	if err != nil || ok {
		t.Fatalf("unknown permission must deny: ok=%v err=%v", ok, err)
	}
}

func TestPrincipalHasRoleCaseInsensitive(t *testing.T) {
	p := &entities.Principal{Roles: []*entities.Role{{Name: "Admin"}}}
	if !p.HasRole("admin") || p.HasRole("superuser") {
		t.Fatal("HasRole should match case-insensitively")
	}
	if !p.HasPermission("x") && len(p.Permissions) == 0 {
		// HasPermission with no permissions must be false
		if p.HasPermission("x") {
			t.Fatal("HasPermission must be false without permissions")
		}
	}
}
