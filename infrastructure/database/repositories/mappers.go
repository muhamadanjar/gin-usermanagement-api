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
		ID:                  m.ID,
		Username:            m.Username,
		Email:               m.Email,
		Name:                m.Name,
		Password:            m.Password,
		FirstName:           m.FirstName,
		LastName:            m.LastName,
		IsVerified:          m.IsVerified,
		IsSuperuser:         m.IsSuperuser,
		IsActive:            m.IsActive,
		Avatar:              m.Avatar,
		EmailVerifiedAt:     m.EmailVerifiedAt,
		FailedLoginAttempts: m.FailedLoginAttempts,
		LockedUntil:         m.LockedUntil,
		LastLogin:           m.LastLogin,
		Status:              m.Status,
		CreatedAt:           m.CreatedAt,
		UpdatedAt:           m.UpdatedAt,
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
		ID:                  e.ID,
		Username:            e.Username,
		Email:               e.Email,
		Name:                e.Name,
		Password:            e.Password,
		FirstName:           e.FirstName,
		LastName:            e.LastName,
		IsVerified:          e.IsVerified,
		IsSuperuser:         e.IsSuperuser,
		IsActive:            e.IsActive,
		Avatar:              e.Avatar,
		EmailVerifiedAt:     e.EmailVerifiedAt,
		FailedLoginAttempts: e.FailedLoginAttempts,
		LockedUntil:         e.LockedUntil,
		LastLogin:           e.LastLogin,
		Status:              e.Status,
		CreatedAt:           e.CreatedAt,
		UpdatedAt:           e.UpdatedAt,
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
		IsActive:    m.IsActive,
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
		IsActive:    e.IsActive,
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
		ID:            m.ID,
		Name:          m.Name,
		Url:           m.Url,
		PermissionKey: m.PermissionKey,
		Icon:          m.Icon,
		Description:   m.Description,
		ParentID:      m.ParentID,
		Sequence:      m.Sequence,
		IsActive:      m.IsActive,
		IsVisible:     m.IsVisible,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
	menu.Parent = toEntityMenu(m.Parent)
	for _, c := range m.Children {
		menu.Children = append(menu.Children, toEntityMenu(c))
	}
	for _, p := range m.Permissions {
		menu.Permissions = append(menu.Permissions, toEntityPermission(p))
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
		ID:            e.ID,
		Name:          e.Name,
		Url:           e.Url,
		PermissionKey: e.PermissionKey,
		Icon:          e.Icon,
		Description:   e.Description,
		ParentID:      e.ParentID,
		Sequence:      e.Sequence,
		IsActive:      e.IsActive,
		IsVisible:     e.IsVisible,
		CreatedAt:     e.CreatedAt,
		UpdatedAt:     e.UpdatedAt,
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

func toEntityTokenHistory(m *models.TokenHistoryModel) *entities.TokenHistory {
	if m == nil {
		return nil
	}
	return &entities.TokenHistory{
		ID:         m.ID,
		UserID:     m.UserID,
		Token:      m.Token,
		ExpiredAt:  m.ExpiredAt,
		LastUsedAt: m.LastUsedAt,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}
}

func toModelTokenHistory(e *entities.TokenHistory) *models.TokenHistoryModel {
	if e == nil {
		return nil
	}
	return &models.TokenHistoryModel{
		ID:         e.ID,
		UserID:     e.UserID,
		Token:      e.Token,
		ExpiredAt:  e.ExpiredAt,
		LastUsedAt: e.LastUsedAt,
	}
}
