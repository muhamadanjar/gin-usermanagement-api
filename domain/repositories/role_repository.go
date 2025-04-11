package repositories

import (
	"usermanagement-api/domain/entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoleRepository interface {
	Create(role *entities.Role) error
	FindByID(id uuid.UUID) (*entities.Role, error)
	FindByName(name string) (*entities.Role, error)
	FindAll(page, pageSize int) ([]*entities.Role, int64, error)
	Update(role *entities.Role) error
	Delete(id uuid.UUID) error
	AssignPermissions(roleID uuid.UUID, permissionIDs []uuid.UUID) error
	FindRolesByUserID(userID uuid.UUID) ([]*entities.Role, error)
	FindPermissionsByRoleIDs(roleIDs []uuid.UUID) ([]*entities.Permission, error)
}

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db}
}

func (r *roleRepository) Create(role *entities.Role) error {
	return r.db.Create(role).Error
}

func (r *roleRepository) FindByID(id uuid.UUID) (*entities.Role, error) {
	var role entities.Role
	if err := r.db.Preload("Permissions").First(&role, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) FindByName(name string) (*entities.Role, error) {
	var role entities.Role
	if err := r.db.Where("name = ?", name).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) FindAll(page, pageSize int) ([]*entities.Role, int64, error) {
	var roles []*entities.Role
	var count int64

	offset := (page - 1) * pageSize

	if err := r.db.Model(&entities.Role{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Preload("Permissions").Offset(offset).Limit(pageSize).Find(&roles).Error; err != nil {
		return nil, 0, err
	}

	return roles, count, nil
}

func (r *roleRepository) Update(role *entities.Role) error {
	return r.db.Save(role).Error
}

func (r *roleRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&entities.Role{}, "id = ?", id).Error
}

func (r *roleRepository) AssignPermissions(roleID uuid.UUID, permissionIDs []uuid.UUID) error {
	tx := r.db.Begin()

	// Remove existing permissions
	if err := tx.Model(&entities.Role{ID: roleID}).Association("Permissions").Clear(); err != nil {
		tx.Rollback()
		return err
	}

	// Add new permissions
	var permissions []*entities.Permission
	for _, permissionID := range permissionIDs {
		permissions = append(permissions, &entities.Permission{ID: permissionID})
	}

	if err := tx.Model(&entities.Role{ID: roleID}).Association("Permissions").Append(permissions); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *roleRepository) FindRolesByUserID(userID uuid.UUID) ([]*entities.Role, error) {
	var roles []*entities.Role
	if err := r.db.Model(&entities.User{ID: userID}).Association("Roles").Find(&roles); err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *roleRepository) FindPermissionsByRoleIDs(roleIDs []uuid.UUID) ([]*entities.Permission, error) {
	var permissions []*entities.Permission
	err := r.db.Table("permissions").
		Joins("INNER JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Where("role_permissions.role_id IN ?", roleIDs).
		Group("permissions.id").
		Find(&permissions).Error

	return permissions, err
}
