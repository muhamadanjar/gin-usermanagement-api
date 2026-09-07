package repositories

import (
	"usermanagement-api/domain/entities"
	"usermanagement-api/infrastructure/database/models"
)

// This file owns the (de)serialization between pure domain entities and GORM
// models. Repositories speak domain entities at the seam; GORM never leaks.

func toEntityUser(m *models.UserModel) *entities.User {
	if m == nil {
		return nil
	}
	u := &entities.User{
		ID:          m.ID,
		Username:    m.Username,
		Email:       m.Email,
		Password:    m.Password,
		FirstName:   m.FirstName,
		LastName:    m.LastName,
		IsSuperuser: m.IsSuperuser,
		IsActive:    m.IsActive,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
	for _, r := range m.Roles {
		u.Roles = append(u.Roles, toEntityRole(r))
	}
	return u
}

func toModelUser(e *entities.User) *models.UserModel {
	if e == nil {
		return nil
	}
	m := &models.UserModel{
		ID:          e.ID,
		Username:    e.Username,
		Email:       e.Email,
		Password:    e.Password,
		FirstName:   e.FirstName,
		LastName:    e.LastName,
		IsSuperuser: e.IsSuperuser,
		IsActive:    e.IsActive,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
	for _, r := range e.Roles {
		m.Roles = append(m.Roles, toModelRole(r))
	}
	return m
}

func toEntityRole(m *models.RoleModel) *entities.Role {
	if m == nil {
		return nil
	}
	r := &entities.Role{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
	for _, p := range m.Permissions {
		r.Permissions = append(r.Permissions, toEntityPermission(p))
	}
	return r
}

func toModelRole(e *entities.Role) *models.RoleModel {
	if e == nil {
		return nil
	}
	m := &models.RoleModel{
		ID:          e.ID,
		Name:        e.Name,
		Description: e.Description,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
	for _, p := range e.Permissions {
		m.Permissions = append(m.Permissions, toModelPermission(p))
	}
	return m
}

func toEntityPermission(m *models.PermissionModel) *entities.Permission {
	if m == nil {
		return nil
	}
	return &entities.Permission{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func toModelPermission(e *entities.Permission) *models.PermissionModel {
	if e == nil {
		return nil
	}
	return &models.PermissionModel{
		ID:          e.ID,
		Name:        e.Name,
		Description: e.Description,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}

func toEntityMenu(m *models.MenuModel) *entities.Menu {
	if m == nil {
		return nil
	}
	menu := &entities.Menu{
		ID:          m.ID,
		Name:        m.Name,
		Url:         m.Url,
		Icon:        m.Icon,
		Description: m.Description,
		ParentID:    m.ParentID,
		Sequence:    m.Sequence,
		IsActive:    m.IsActive,
		IsVisible:   m.IsVisible,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
	menu.Parent = toEntityMenu(m.Parent)
	for _, c := range m.Children {
		menu.Children = append(menu.Children, toEntityMenu(c))
	}
	return menu
}

// toModelMenu maps the flat, writable fields only — parent/children are read
// back via Preload by the repository, never written through this struct.
func toModelMenu(e *entities.Menu) *models.MenuModel {
	if e == nil {
		return nil
	}
	return &models.MenuModel{
		ID:          e.ID,
		Name:        e.Name,
		Url:         e.Url,
		Icon:        e.Icon,
		Description: e.Description,
		ParentID:    e.ParentID,
		Sequence:    e.Sequence,
		IsActive:    e.IsActive,
		IsVisible:   e.IsVisible,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}

func toEntityModelPermission(m *models.ModelPermissionModel) *entities.ModelPermission {
	if m == nil {
		return nil
	}
	mp := &entities.ModelPermission{
		ID:           m.ID,
		ModelID:      m.ModelID,
		ModelType:    m.ModelType,
		PermissionID: m.PermissionID,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
	mp.Permission = *toEntityPermission(&m.Permission)
	return mp
}

func toModelModelPermission(e *entities.ModelPermission) *models.ModelPermissionModel {
	if e == nil {
		return nil
	}
	return &models.ModelPermissionModel{
		ID:           e.ID,
		ModelID:      e.ModelID,
		ModelType:    e.ModelType,
		PermissionID: e.PermissionID,
		CreatedAt:    e.CreatedAt,
		UpdatedAt:    e.UpdatedAt,
	}
}

func toEntityUserMeta(m *models.UserMetaModel) *entities.UserMeta {
	if m == nil {
		return nil
	}
	return &entities.UserMeta{
		ID:     m.ID,
		Key:    m.Key,
		Value:  m.Value,
		UserID: m.UserID,
	}
}

func toModelUserMeta(e *entities.UserMeta) *models.UserMetaModel {
	if e == nil {
		return nil
	}
	return &models.UserMetaModel{
		ID:     e.ID,
		Key:    e.Key,
		Value:  e.Value,
		UserID: e.UserID,
	}
}

func toEntitySetting(m *models.SettingModel) *entities.Setting {
	if m == nil {
		return nil
	}
	return &entities.Setting{Key: m.Key, Value: m.Value}
}

func toModelSetting(e *entities.Setting) *models.SettingModel {
	if e == nil {
		return nil
	}
	return &models.SettingModel{Key: e.Key, Value: e.Value}
}
