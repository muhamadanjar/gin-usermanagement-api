// internal/dto/setting_dto.go
package dto

type CreateSettingRequest struct {
	Key   string `json:"key" binding:"required"`
	Value string `json:"value" binding:"required"`
}

type UpdateSettingRequest struct {
	Value string `json:"value" binding:"required"`
}

type SettingResponse struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
