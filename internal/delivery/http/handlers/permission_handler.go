package handlers

import (
	"net/http"
	"strconv"
	"usermanagement-api/internal/dto"
	"usermanagement-api/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PermissionHandler struct {
	permissionUseCase usecase.PermissionUseCase
}

func NewPermissionHandler(permissionUseCase usecase.PermissionUseCase) *PermissionHandler {
	return &PermissionHandler{
		permissionUseCase: permissionUseCase,
	}
}

func (h *PermissionHandler) CreatePermission(c *gin.Context) {
	var req dto.CreatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.permissionUseCase.Create(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// GetPermission godoc
// @Summary Get permission
// @Description Get permission by ID
// @Tags permissions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Permission ID"
// @Success 200 {object} dto.PermissionResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /permissions/{id} [get]
func (h *PermissionHandler) GetPermission(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid permission id"})
		return
	}

	resp, err := h.permissionUseCase.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "permission not found"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetAllPermissions godoc
// @Summary Get all permissions
// @Description Get all permissions with pagination
// @Tags permissions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 10)"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Router /permissions [get]
func (h *PermissionHandler) GetAllPermissions(c *gin.Context) {
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

	permissions, total, err := h.permissionUseCase.GetAll(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": permissions,
		"meta": gin.H{
			"page":       page,
			"page_size":  pageSize,
			"total":      total,
			"total_page": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// UpdatePermission godoc
// @Summary Update permission
// @Description Update permission by ID
// @Tags permissions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Permission ID"
// @Param permission body dto.UpdatePermissionRequest true "Permission information"
// @Success 200 {object} dto.PermissionResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /permissions/{id} [put]
func (h *PermissionHandler) UpdatePermission(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid permission id"})
		return
	}

	var req dto.UpdatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.permissionUseCase.Update(id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DeletePermission godoc
// @Summary Delete permission
// @Description Delete permission by ID
// @Tags permissions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Permission ID"
// @Success 204 {object} nil
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /permissions/{id} [delete]
func (h *PermissionHandler) DeletePermission(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid permission id"})
		return
	}

	if err := h.permissionUseCase.Delete(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
