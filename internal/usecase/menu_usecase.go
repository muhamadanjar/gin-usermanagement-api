package usecase

import (
	"errors"
	"time"
	"usermanagement-api/domain/entities"
	"usermanagement-api/domain/repositories"
	"usermanagement-api/internal/dto"

	"github.com/google/uuid"
)

type MenuUseCase interface {
	Create(req *dto.CreateMenuRequest) (*dto.MenuResponse, error)
	GetByID(id uuid.UUID) (*dto.MenuResponse, error)
	GetAll(page, pageSize int) ([]*dto.MenuResponse, int64, error)
	GetAllActive() ([]*dto.MenuResponse, error)
	Update(id uuid.UUID, req *dto.UpdateMenuRequest) (*dto.MenuResponse, error)
	Delete(id uuid.UUID) error
}

type menuUseCase struct {
	menuRepo repositories.MenuRepository
}

func NewMenuUseCase(menuRepo repositories.MenuRepository) MenuUseCase {
	return &menuUseCase{
		menuRepo: menuRepo,
	}
}

func (uc *menuUseCase) Create(req *dto.CreateMenuRequest) (*dto.MenuResponse, error) {
	// Check if menu name already exists
	if _, err := uc.menuRepo.FindByName(req.Name); err == nil {
		return nil, errors.New("menu name already exists")
	}

	// Create menu
	menu := &entities.Menu{
		Name:        req.Name,
		Url:         req.Url,
		Icon:        req.Icon,
		Description: req.Description,
		ParentID:    req.ParentID,
		Sequence:    req.Sequence,
		IsActive:    true,
	}

	// Save menu
	if err := uc.menuRepo.Create(menu); err != nil {
		return nil, err
	}

	// Reload to get parent/children
	loadedMenu, err := uc.menuRepo.FindByID(menu.ID)
	if err != nil {
		return nil, err
	}

	return uc.mapToMenuResponse(loadedMenu), nil
}

func (uc *menuUseCase) GetByID(id uuid.UUID) (*dto.MenuResponse, error) {
	menu, err := uc.menuRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return uc.mapToMenuResponse(menu), nil
}

func (uc *menuUseCase) GetAll(page, pageSize int) ([]*dto.MenuResponse, int64, error) {
	menus, total, err := uc.menuRepo.FindAll(page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	var response []*dto.MenuResponse
	for _, menu := range menus {
		response = append(response, uc.mapToMenuSimpleResponse(menu))
	}

	return response, total, nil
}

func (uc *menuUseCase) GetAllActive() ([]*dto.MenuResponse, error) {
	menus, err := uc.menuRepo.FindAllActive()
	if err != nil {
		return nil, err
	}

	var response []*dto.MenuResponse
	for _, menu := range menus {
		response = append(response, uc.mapToMenuSimpleResponse(menu))
	}

	return response, nil
}

func (uc *menuUseCase) Update(id uuid.UUID, req *dto.UpdateMenuRequest) (*dto.MenuResponse, error) {
	menu, err := uc.menuRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Name != "" && req.Name != menu.Name {
		// Check if new name already exists
		if existingMenu, err := uc.menuRepo.FindByName(req.Name); err == nil && existingMenu.ID != id {
			return nil, errors.New("menu name already exists")
		}
		menu.Name = req.Name
	}

	if req.Url != "" {
		menu.Url = req.Url
	}

	if req.Icon != "" {
		menu.Icon = req.Icon
	}

	if req.Description != "" {
		menu.Description = req.Description
	}

	if req.ParentID != nil {
		menu.ParentID = req.ParentID
	}

	if req.Sequence != 0 {
		menu.Sequence = req.Sequence
	}

	if req.IsActive != nil {
		menu.IsActive = *req.IsActive
	}

	// Update menu
	if err := uc.menuRepo.Update(menu); err != nil {
		return nil, err
	}

	// Reload to get updated parent/children
	loadedMenu, err := uc.menuRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return uc.mapToMenuResponse(loadedMenu), nil
}

func (uc *menuUseCase) Delete(id uuid.UUID) error {
	menu, err := uc.menuRepo.FindByID(id)
	if err != nil {
		return err
	}

	// Check if menu has children
	if len(menu.Children) > 0 {
		return errors.New("cannot delete menu with children")
	}

	// Delete menu
	return uc.menuRepo.Delete(id)
}

func (uc *menuUseCase) mapToMenuSimpleResponse(menu *entities.Menu) *dto.MenuResponse {
	resp := &dto.MenuResponse{
		ID:          menu.ID,
		Name:        menu.Name,
		Url:         menu.Url,
		Icon:        menu.Icon,
		Description: menu.Description,
		ParentID:    menu.ParentID,
		Sequence:    menu.Sequence,
		IsActive:    menu.IsActive,
		CreatedAt:   menu.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   menu.UpdatedAt.Format(time.RFC3339),
	}
	return resp

}

func (uc *menuUseCase) mapToMenuResponse(menu *entities.Menu) *dto.MenuResponse {
	resp := &dto.MenuResponse{
		ID:          menu.ID,
		Name:        menu.Name,
		Url:         menu.Url,
		Icon:        menu.Icon,
		Description: menu.Description,
		ParentID:    menu.ParentID,
		Sequence:    menu.Sequence,
		IsActive:    menu.IsActive,
		CreatedAt:   menu.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   menu.UpdatedAt.Format(time.RFC3339),
	}

	// Map parent if exists
	if menu.Parent != nil {
		resp.Parent = &dto.MenuSimple{
			ID:   menu.Parent.ID,
			Name: menu.Parent.Name,
			Url:  menu.Parent.Url,
		}
	}

	// Map children if exists
	if menu.Children != nil {
		for _, child := range menu.Children {
			resp.Children = append(resp.Children, &dto.MenuSimple{
				ID:   child.ID,
				Name: child.Name,
				Url:  child.Url,
			})
		}
	}

	return resp
}
