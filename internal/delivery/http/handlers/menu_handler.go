package handlers

import (
	"net/http"
	"strconv"
	"usermanagement-api/internal/dto"
	"usermanagement-api/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MenuHandler struct {
	menuUseCase usecase.MenuUseCase
}

func NewMenuHandler(menuUseCase usecase.MenuUseCase) *MenuHandler {
	return &MenuHandler{
		menuUseCase: menuUseCase,
	}
}

// CreateMenu godoc
// @Summary Create menu
// @Description Create a new menu
// @Tags menus
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param menu body dto.CreateMenuRequest true "Menu information"
// @Success 201 {object} dto.MenuResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /menus [post]
func (h *MenuHandler) CreateMenu(c *gin.Context) {
	var req dto.CreateMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.menuUseCase.Create(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// GetMenu godoc
// @Summary Get menu
// @Description Get menu by ID
// @Tags menus
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Menu ID"
// @Success 200 {object} dto.MenuResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /menus/{id} [get]
func (h *MenuHandler) GetMenu(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid menu id"})
		return
	}

	resp, err := h.menuUseCase.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "menu not found"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetAllMenus godoc
// @Summary Get all menus
// @Description Get all menus with pagination
// @Tags menus
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 10)"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Router /menus [get]
func (h *MenuHandler) GetAllMenus(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	menus, total, err := h.menuUseCase.GetAll(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"menus": menus,
		"meta": gin.H{
			"page":       page,
			"page_size":  pageSize,
			"total":      total,
			"total_page": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// GetActiveMenus godoc
// @Summary Get active menus
// @Description Get all active menus
// @Tags menus
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.MenuResponse
// @Failure 401 {object} map[string]string
// @Router /menus/active [get]
func (h *MenuHandler) GetActiveMenus(c *gin.Context) {
	menus, err := h.menuUseCase.GetAllActive()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, menus)
}

// UpdateMenu godoc
// @Summary Update menu
// @Description Update menu by ID
// @Tags menus
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Menu ID"
// @Param menu body dto.UpdateMenuRequest true "Menu information"
// @Success 200 {object} dto.MenuResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /menus/{id} [put]
func (h *MenuHandler) UpdateMenu(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid menu id"})
		return
	}

	var req dto.UpdateMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.menuUseCase.Update(id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteMenu godoc
// @Summary Delete menu
// @Description Delete menu by ID
// @Tags menus
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Menu ID"
// @Success 204 {object} nil
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /menus/{id} [delete]
func (h *MenuHandler) DeleteMenu(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid menu id"})
		return
	}

	if err := h.menuUseCase.Delete(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
