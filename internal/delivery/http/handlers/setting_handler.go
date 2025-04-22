// internal/delivery/http/handlers/setting_handler.go
package handlers

import (
	"net/http"
	"usermanagement-api/internal/dto"
	"usermanagement-api/internal/usecase"

	"github.com/gin-gonic/gin"
)

type SettingHandler struct {
	settingUseCase usecase.SettingUseCase
}

func NewSettingHandler(settingUseCase usecase.SettingUseCase) *SettingHandler {
	return &SettingHandler{
		settingUseCase: settingUseCase,
	}
}

func (h *SettingHandler) CreateOrUpdate(c *gin.Context) {
	var req dto.CreateSettingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.settingUseCase.CreateOrUpdate(req.Key, req.Value); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Setting updated successfully"})
}

func (h *SettingHandler) GetByKey(c *gin.Context) {
	key := c.Param("key")

	resp, err := h.settingUseCase.GetByKey(key)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "setting not found"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *SettingHandler) GetAll(c *gin.Context) {
	settingsMap, err := h.settingUseCase.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, settingsMap)
}

func (h *SettingHandler) Delete(c *gin.Context) {
	key := c.Param("key")

	if err := h.settingUseCase.Delete(key); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
