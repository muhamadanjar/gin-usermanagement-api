package repositories

import (
	"usermanagement-api/domain/entities"
	domainrepos "usermanagement-api/domain/repositories"
	"usermanagement-api/infrastructure/database/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) domainrepos.RoleRepository {
	return &roleRepository{db}
}

func (r *roleRepository) Create(role *entities.Role) error {
	m := toModelRole(role)
	if err := r.db.Create(m).Error; err != nil {
		return err
	}
	role.ID = m.ID
	role.CreatedAt = m.CreatedAt
	role.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *roleRepository) FindByID(id uuid.UUID) (*entities.Role, error) {
	var m models.RoleModel
	if err := r.db.Preload("Permissions").First(&m, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return toEntityRole(&m), nil
}

func (r *roleRepository) FindByName(name string) (*entities.Role, error) {
	var m models.RoleModel
	if err := r.db.Where("name = ?", name).First(&m).Error; err != nil {
		return nil, err
	}
	return toEntityRole(&m), nil
}

func (r *roleRepository) FindAll(page, pageSize int) ([]*entities.Role, int64, error) {
	var rows []*models.RoleModel
	var count int64

	offset := (page - 1) * pageSize

	if err := r.db.Model(&models.RoleModel{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Preload("Permissions").Offset(offset).Limit(pageSize).Find(&rows).Error; err != nil {
		return nil, 0, err
	}

	roles := make([]*entities.Role, 0, len(rows))
	for _, m := range rows {
		roles = append(roles, toEntityRole(m))
	}
	return roles, count, nil
}

func (r *roleRepository) Update(role *entities.Role) error {
	return r.db.Save(toModelRole(role)).Error
}

func (r *roleRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.RoleModel{}, "id = ?", id).Error
}

func (r *roleRepository) AssignPermissions(roleID uuid.UUID, permissionIDs []uuid.UUID) error {
	tx := r.db.Begin()

	if err := tx.Model(&models.RoleModel{ID: roleID}).Association("Permissions").Clear(); err != nil {
		tx.Rollback()
		return err
	}

	var permissions []*models.PermissionModel
	for _, permissionID := range permissionIDs {
		permissions = append(permissions, &models.PermissionModel{ID: permissionID})
	}

	if err := tx.Model(&models.RoleModel{ID: roleID}).Association("Permissions").Append(permissions); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *roleRepository) FindRolesByUserID(userID uuid.UUID) ([]*entities.Role, error) {
	var rows []*models.RoleModel
	if err := r.db.Model(&models.UserModel{ID: userID}).Association("Roles").Find(&rows); err != nil {
		return nil, err
	}
	roles := make([]*entities.Role, 0, len(rows))
	for _, m := range rows {
		roles = append(roles, toEntityRole(m))
	}
	return roles, nil
}

func (r *roleRepository) FindPermissionsByRoleIDs(roleIDs []uuid.UUID) ([]*entities.Permission, error) {
	var rows []*models.PermissionModel
	err := r.db.Table("permissions").
		Joins("INNER JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Where("role_permissions.role_id IN ?", roleIDs).
		Group("permissions.id").
		Find(&rows).Error

	permissions := make([]*entities.Permission, 0, len(rows))
	for _, m := range rows {
		permissions = append(permissions, toEntityPermission(m))
	}
	return permissions, err
}
